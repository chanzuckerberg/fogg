package apply

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/afero"
)

var foggManagedRoots = []string{
	"terraform", ".github", "scripts", ".circleci",
	".fogg-version", "Makefile", ".gitignore", ".gitattributes",
	".terraformignore", ".travis.yml", "README.md",
	".terraform.d", "terraform.d",
}

func ApplyDryRun(fs afero.Fs, repoRoot string, conf *v2.Config, tmpl *templates.T, upgrade bool) (diff string, hasChanges bool, err error) {
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

	diffOutput, hasChanges, err := diffFoggManagedPaths(repoRoot, tempDir)
	if err != nil {
		return "", false, errs.WrapUser(err, "unable to compute diff")
	}

	return diffOutput, hasChanges, nil
}

var copySkipPrefixes = []string{
	".terraform.d/plugin-cache",
	".terraform.d/versions",
}

func shouldSkipCopy(path string) bool {
	path = filepath.ToSlash(path)
	for _, prefix := range copySkipPrefixes {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	if strings.Contains(path, "/.terraform/") {
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

func diffFoggManagedPaths(currentDir, plannedDir string) (string, bool, error) {
	currentFiles := make(map[string]struct{})
	plannedFiles := make(map[string]struct{})

	for _, root := range foggManagedRoots {
		if err := collectPaths(currentDir, root, currentFiles); err != nil {
			return "", false, err
		}
		if err := collectPathsFromOS(plannedDir, root, plannedFiles); err != nil {
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

	dmp := diffmatchpatch.New()
	var sb strings.Builder
	hasChanges := false

	for _, relPath := range paths {
		inCurrent := inSet(currentFiles, relPath)
		inPlanned := inSet(plannedFiles, relPath)

		var currentContent, plannedContent string
		if inCurrent {
			data, err := os.ReadFile(filepath.Join(currentDir, relPath))
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
			sb.WriteString(formatUnifiedDiff(dmp, "", plannedContent))
			sb.WriteString("\n")
		} else if inCurrent && !inPlanned {
			hasChanges = true
			fmt.Fprintf(&sb, "--- %s\n+++ /dev/null\n", relPath)
			sb.WriteString(formatUnifiedDiff(dmp, currentContent, ""))
			sb.WriteString("\n")
		} else if inCurrent && inPlanned && currentContent != plannedContent {
			hasChanges = true
			fmt.Fprintf(&sb, "--- %s\n+++ %s\n", relPath, relPath)
			sb.WriteString(formatUnifiedDiff(dmp, currentContent, plannedContent))
			sb.WriteString("\n")
		}
	}

	return sb.String(), hasChanges, nil
}

func inSet(m map[string]struct{}, key string) bool {
	_, ok := m[key]
	return ok
}

func collectPaths(baseDir, root string, out map[string]struct{}) error {
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

func collectPathsFromOS(baseDir, root string, out map[string]struct{}) error {
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

func formatUnifiedDiff(dmp *diffmatchpatch.DiffMatchPatch, oldText, newText string) string {
	diffs := dmp.DiffMain(oldText, newText, true)
	dmp.DiffCleanupSemantic(diffs)

	var sb strings.Builder
	oldCount, newCount := 0, 0

	for _, d := range diffs {
		lines := strings.Split(d.Text, "\n")
		if len(lines) > 1 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		for _, line := range lines {
			switch d.Type {
			case diffmatchpatch.DiffDelete:
				fmt.Fprintf(&sb, "-%s\n", line)
				oldCount++
			case diffmatchpatch.DiffInsert:
				fmt.Fprintf(&sb, "+%s\n", line)
				newCount++
			case diffmatchpatch.DiffEqual:
				fmt.Fprintf(&sb, " %s\n", line)
				oldCount++
				newCount++
			}
		}
	}

	if sb.Len() == 0 {
		return ""
	}

	var header string
	if oldCount == 0 {
		header = fmt.Sprintf("@@ -0,0 +1,%d @@\n", newCount)
	} else if newCount == 0 {
		header = fmt.Sprintf("@@ -1,%d +0,0 @@\n", oldCount)
	} else {
		header = fmt.Sprintf("@@ -1,%d +1,%d @@\n", oldCount, newCount)
	}
	return header + sb.String()
}
