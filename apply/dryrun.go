package apply

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/spf13/afero"
)

var foggManagedRoots = []string{
	"terraform", ".github", "scripts", ".circleci",
	".fogg-version", "Makefile", ".gitignore", ".gitattributes",
	".terraformignore", ".travis.yml", "README.md",
	".terraform.d", "terraform.d",
}

func ApplyDryRun(fs afero.Fs, conf *v2.Config, tmpl *templates.T, upgrade bool) (diff string, hasChanges bool, err error) {
	tempDir, err := os.MkdirTemp("", "fogg-dryrun-")
	if err != nil {
		return "", false, errs.WrapInternal(err, "unable to create temp dir for dry-run")
	}
	defer os.RemoveAll(tempDir)

	if err := copyFoggManagedPaths(fs, tempDir); err != nil {
		return "", false, errs.WrapUser(err, "unable to copy current state for dry-run")
	}

	tempFs := afero.NewBasePathFs(afero.NewOsFs(), tempDir)
	if err := Apply(tempFs, conf, tmpl, upgrade); err != nil {
		return "", false, err
	}

	diffOutput, hasChanges, err := diffFoggManagedPaths(fs, tempDir)
	if err != nil {
		return "", false, errs.WrapUser(err, "unable to compute diff")
	}

	return diffOutput, hasChanges, nil
}

var copySkipPrefixes = []string{
	".git",
	".terraform.d/versions",
}

func shouldSkipCopy(path string) bool {
	path = filepath.ToSlash(path)
	if path == ".terraform.d/plugin-cache/.gitignore" {
		return false
	}
	if path == ".terraform.d/plugin-cache" || strings.HasPrefix(path, ".terraform.d/plugin-cache/") {
		return true
	}
	for _, prefix := range copySkipPrefixes {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	if strings.Contains(path, "/.terraform/") {
		return true
	}
	if strings.HasPrefix(path, "terraform.d/plugins") || strings.HasPrefix(path, "terraform.d/modules") {
		return true
	}
	return false
}

func copyFoggManagedPaths(src afero.Fs, destDir string) error {
	for _, root := range foggManagedRoots {
		info, err := src.Stat(root)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		destPath := filepath.Join(destDir, root)
		if info.IsDir() {
			if err := copyDir(src, root, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(src, root, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyDir(src afero.Fs, srcPath, destPath string) error {
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}
	return afero.Walk(src, srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if shouldSkipCopy(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		rel, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destPath, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(src, path, target)
	})
}

func copyFile(src afero.Fs, srcPath, destPath string) error {
	data, err := afero.ReadFile(src, srcPath)
	if err != nil {
		if errors.Is(err, syscall.EISDIR) {
			return nil
		}
		return err
	}
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(destPath, data, 0644)
}

func diffFoggManagedPaths(currentFs afero.Fs, plannedDir string) (string, bool, error) {
	currentFiles := make(map[string]struct{})
	plannedFiles := make(map[string]struct{})

	for _, root := range foggManagedRoots {
		if err := collectPathsFromFs(currentFs, root, currentFiles); err != nil {
			return "", false, err
		}
		if err := collectPathsOS(plannedDir, root, plannedFiles); err != nil {
			return "", false, err
		}
	}

	allPaths := make(map[string]struct{})
	for p := range currentFiles {
		allPaths[p] = struct{}{}
	}
	for p := range plannedFiles {
		allPaths[p] = struct{}{}
	}

	var paths []string
	for p := range allPaths {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	var sb strings.Builder
	hasChanges := false

	for _, relPath := range paths {
		inCurrent := inSet(currentFiles, relPath)
		inPlanned := inSet(plannedFiles, relPath)

		var currentContent, plannedContent string
		if inCurrent {
			data, err := afero.ReadFile(currentFs, relPath)
			if err != nil {
				if errors.Is(err, syscall.EISDIR) {
					continue
				}
				return "", false, err
			}
			currentContent = string(data)
		}
		if inPlanned {
			data, err := os.ReadFile(filepath.Join(plannedDir, relPath))
			if err != nil {
				if errors.Is(err, syscall.EISDIR) {
					continue
				}
				return "", false, err
			}
			plannedContent = string(data)
		}

		if !inCurrent && inPlanned {
			hasChanges = true
			fmt.Fprintf(&sb, "--- /dev/null\n+++ %s\n", relPath)
			sb.WriteString(formatUnifiedDiff("", plannedContent))
			sb.WriteString("\n")
		} else if inCurrent && !inPlanned {
			hasChanges = true
			fmt.Fprintf(&sb, "--- %s\n+++ /dev/null\n", relPath)
			sb.WriteString(formatUnifiedDiff(currentContent, ""))
			sb.WriteString("\n")
		} else if inCurrent && inPlanned && currentContent != plannedContent {
			hasChanges = true
			fmt.Fprintf(&sb, "--- %s\n+++ %s\n", relPath, relPath)
			sb.WriteString(formatUnifiedDiff(currentContent, plannedContent))
			sb.WriteString("\n")
		}
	}

	return sb.String(), hasChanges, nil
}

func inSet(m map[string]struct{}, key string) bool {
	_, ok := m[key]
	return ok
}

func collectPathsFromFs(fs afero.Fs, root string, out map[string]struct{}) error {
	info, err := fs.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		out[root] = struct{}{}
		return nil
	}
	return afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if shouldSkipCopy(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		if !info.IsDir() {
			out[path] = struct{}{}
		}
		return nil
	})
}

func collectPathsOS(baseDir, root string, out map[string]struct{}) error {
	fullPath := filepath.Join(baseDir, root)
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		out[root] = struct{}{}
		return nil
	}
	return filepath.WalkDir(fullPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}
		if shouldSkipCopy(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			out[rel] = struct{}{}
		}
		return nil
	})
}

func formatUnifiedDiff(oldText, newText string) string {
	oldFile, err := os.CreateTemp("", "fogg-diff-old-")
	if err != nil {
		return fallbackUnifiedDiff(oldText, newText)
	}
	oldPath := oldFile.Name()
	defer os.Remove(oldPath)
	defer func() {
		if err := oldFile.Close(); err != nil {
			// best-effort close on defer; error cannot propagate
		}
	}()

	newFile, err := os.CreateTemp("", "fogg-diff-new-")
	if err != nil {
		return fallbackUnifiedDiff(oldText, newText)
	}
	newPath := newFile.Name()
	defer os.Remove(newPath)
	defer func() {
		if err := newFile.Close(); err != nil {
			// best-effort close on defer; error cannot propagate
		}
	}()

	if _, err := oldFile.WriteString(oldText); err != nil {
		return fallbackUnifiedDiff(oldText, newText)
	}
	if _, err := newFile.WriteString(newText); err != nil {
		return fallbackUnifiedDiff(oldText, newText)
	}
	if err := oldFile.Close(); err != nil {
		return fallbackUnifiedDiff(oldText, newText)
	}
	if err := newFile.Close(); err != nil {
		return fallbackUnifiedDiff(oldText, newText)
	}

	cmd := exec.Command("git", "diff", "--no-index", "--no-color", oldPath, newPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
			return fallbackUnifiedDiff(oldText, newText)
		}
	}

	diff := string(out)
	if idx := strings.Index(diff, "@@"); idx >= 0 {
		return diff[idx:]
	}
	return fallbackUnifiedDiff(oldText, newText)
}

func fallbackUnifiedDiff(oldText, newText string) string {
	oldLines := strings.Split(oldText, "\n")
	newLines := strings.Split(newText, "\n")
	var sb strings.Builder
	for _, line := range oldLines {
		if line != "" || len(oldLines) > 1 {
			fmt.Fprintf(&sb, "-%s\n", line)
		}
	}
	for _, line := range newLines {
		if line != "" || len(newLines) > 1 {
			fmt.Fprintf(&sb, "+%s\n", line)
		}
	}
	if sb.Len() == 0 {
		return ""
	}
	return fmt.Sprintf("@@ -1,%d +1,%d @@\n", len(oldLines), len(newLines)) + sb.String()
}
