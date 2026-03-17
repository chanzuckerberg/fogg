package apply

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// DryRun computes a unified diff of what Apply would change without modifying the filesystem.
func DryRun(fs afero.Fs, conf *v2.Config, tmpl *templates.T, upgrade bool) (diff string, hasChanges bool, err error) {
	tempBase, err := os.MkdirTemp("", "fogg-dryrun-")
	if err != nil {
		return "", false, errs.WrapInternal(err, "unable to create temp dir for dry-run")
	}
	defer os.RemoveAll(tempBase)

	currentDir := filepath.Join(tempBase, "current")
	plannedDir := filepath.Join(tempBase, "planned")
	if err := os.MkdirAll(currentDir, 0755); err != nil {
		return "", false, errs.WrapInternal(err, "unable to create current dir")
	}
	if err := os.MkdirAll(plannedDir, 0755); err != nil {
		return "", false, errs.WrapInternal(err, "unable to create planned dir")
	}

	if err := copyFoggManagedPaths(fs, currentDir); err != nil {
		return "", false, errs.WrapUser(err, "unable to copy current state for dry-run")
	}
	if err := copyFoggManagedPaths(fs, plannedDir); err != nil {
		return "", false, errs.WrapUser(err, "unable to copy current state for dry-run")
	}

	tempFs := afero.NewBasePathFs(afero.NewOsFs(), plannedDir)
	if err := Apply(tempFs, conf, tmpl, upgrade); err != nil {
		return "", false, err
	}

	diffOutput, hasChanges, err := diffTrees(tempBase, "current", "planned")
	if err != nil {
		return "", false, errs.WrapUser(err, "unable to compute diff")
	}

	return diffOutput, hasChanges, nil
}

func diffTrees(baseDir, currentName, plannedName string) (string, bool, error) {
	cmd := exec.Command("git", "diff", "--no-index", "--no-color", currentName, plannedName)
	cmd.Dir = baseDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
			return "", false, err
		}
	}
	diff := string(out)
	hasChanges := strings.Contains(diff, "diff --git")
	if !hasChanges {
		return "", false, nil
	}
	filtered := filterGitDiffToFoggPaths(diff, currentName+"/", plannedName+"/")
	return filtered, true, nil
}

func filterGitDiffToFoggPaths(fullDiff, currentPrefix, plannedPrefix string) string {
	var sb strings.Builder
	lines := strings.Split(fullDiff, "\n")
	i := 0
	for i < len(lines) {
		line := lines[i]
		if strings.HasPrefix(line, "diff --git ") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				aPath := parts[2]
				bPath := parts[3]
				var relPath string
				if strings.HasPrefix(aPath, "a/"+currentPrefix) {
					relPath = strings.TrimPrefix(aPath, "a/"+currentPrefix)
				} else if strings.HasPrefix(bPath, "b/"+plannedPrefix) {
					relPath = strings.TrimPrefix(bPath, "b/"+plannedPrefix)
				}
				if relPath != "" && isFoggManagedPath(relPath) {
					i++
					for i < len(lines) && !strings.HasPrefix(lines[i], "diff --git ") {
						l := lines[i]
						if strings.HasPrefix(l, "--- ") {
							if strings.TrimSpace(strings.TrimPrefix(l, "--- ")) == "/dev/null" {
								sb.WriteString("--- /dev/null\n")
							} else {
								fmt.Fprintf(&sb, "--- %s\n", relPath)
							}
						} else if strings.HasPrefix(l, "+++ ") {
							if strings.TrimSpace(strings.TrimPrefix(l, "+++ ")) == "/dev/null" {
								sb.WriteString("+++ /dev/null\n")
							} else {
								fmt.Fprintf(&sb, "+++ %s\n", relPath)
							}
						} else if strings.HasPrefix(l, "@@") || strings.HasPrefix(l, "-") || strings.HasPrefix(l, "+") || strings.HasPrefix(l, " ") {
							sb.WriteString(l)
							sb.WriteString("\n")
						}
						i++
					}
					sb.WriteString("\n")
					continue
				}
			}
		}
		i++
	}
	return strings.TrimSuffix(sb.String(), "\n\n")
}

func isFoggManagedPath(relPath string) bool {
	relPath = filepath.ToSlash(relPath)
	if strings.HasSuffix(relPath, "/terraform.d") && strings.HasPrefix(relPath, "terraform/") {
		return false
	}
	for _, root := range foggManagedRoots {
		if relPath == root || strings.HasPrefix(relPath, root+"/") {
			return !shouldSkipCopy(relPath)
		}
	}
	return false
}

var copySkipPrefixes = []string{
	".git",
	".terraform.d/versions",
}

func shouldSkipCopy(path string) bool {
	path = filepath.ToSlash(path)
	if path == ".terraform.d/plugin-cache" || path == ".terraform.d/plugin-cache/.gitignore" {
		return false
	}
	if strings.HasPrefix(path, ".terraform.d/plugin-cache/") {
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
		rel, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destPath, rel)
		if info.Mode()&os.ModeSymlink != 0 {
			return copySymlink(src, path, target)
		}
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(src, path, target)
	})
}

func copySymlink(src afero.Fs, srcPath, destPath string) error {
	realPath, ok := realPathForFs(src, srcPath)
	if !ok {
		return nil
	}
	linkTarget, err := os.Readlink(realPath)
	if err != nil {
		return err
	}
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.Symlink(linkTarget, destPath)
}

func realPathForFs(fs afero.Fs, path string) (string, bool) {
	bp, ok := fs.(*afero.BasePathFs)
	if !ok {
		return "", false
	}
	return afero.FullBaseFsPath(bp, path), true
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
