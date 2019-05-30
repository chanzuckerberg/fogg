package plugins

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// TypePluginFormat is the plugin format such as binary, zip, tar
type TypePluginFormat string

const (
	// TypePluginFormatTar is a tar archived plugin
	TypePluginFormatTar TypePluginFormat = "tar"
	// TypePluginFormatBin is a binary plugin
	TypePluginFormatBin TypePluginFormat = "bin"
	// TypePluginFormatZip is a zip archive plugin
	TypePluginFormatZip TypePluginFormat = "zip"
)

// CustomPlugin is a custom plugin
type CustomPlugin struct {
	URL       string           `json:"url" validate:"required"`
	Format    TypePluginFormat `json:"format" validate:"required"`
	TarConfig *TarConfig       `json:"tar_config,omitempty"`
	TargetDir string           `json:"target_dir,omitempty"`
}

// TarConfig configures the tar unpacking
type TarConfig struct {
	StripComponents int `json:"strip_components,omitempty"`
}

func (tc *TarConfig) getStripComponents() int {
	if tc == nil {
		return 0
	}
	return tc.StripComponents
}

// Install installs the custom plugin
func (cp *CustomPlugin) Install(fs afero.Fs, pluginName string) error {
	return cp.install(fs, pluginName, cp.URL, runtime.GOOS, runtime.GOARCH)
}

// Install delegates to install
func (cp *CustomPlugin) install(
	fs afero.Fs,
	pluginName string,
	url string,
	pluginOS string,
	pluginArch string,
) error {
	if cp == nil {
		return errs.NewUser("nil CustomPlugin")
	}
	if fs == nil {
		return errs.NewUser("nil fs")
	}

	fullUrl, err := Template(url, pluginOS, pluginArch)
	if err != nil {
		return err
	}
	tmpPath, err := cp.fetch(pluginName, fullUrl)
	defer os.Remove(tmpPath)
	if err != nil {
		return err
	}
	return cp.process(fs, pluginName, tmpPath, cp.TargetDir)
}

func Template(url, os, arch string) (string, error) {
	data := struct {
		OS   string
		Arch string
	}{os, arch}

	t, err := template.New("url").Parse(url)
	if err != nil {
		return "", errs.WrapUserf(err, "could not parse url template %s", url)
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	out := buf.String()
	return out, nil
}

// SetTargetPath sets the target path for this plugin
func (cp *CustomPlugin) SetTargetPath(path string) {
	cp.TargetDir = path
}

// fetch fetches the custom plugin at URL
func (cp *CustomPlugin) fetch(pluginName string, url string) (string, error) {
	tmpFile, err := ioutil.TempFile("", fmt.Sprintf("%s-*.tmp", pluginName))
	if err != nil {
		return "", errs.WrapUser(err, "could not create temporary directory") //FIXME
	}
	logrus.Debugf("downloading %s to tempfile", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", errs.WrapUserf(err, "could not get %s", url) // FIXME
	}
	defer resp.Body.Close()
	_, err = io.Copy(tmpFile, resp.Body)
	return tmpFile.Name(), errs.WrapUser(err, "could not download file") //FIXME
}

// process the custom plugin
func (cp *CustomPlugin) process(fs afero.Fs, pluginName string, path string, targetDir string) error {
	switch cp.Format {
	case TypePluginFormatTar:
		return cp.processTar(fs, path, targetDir)
	case TypePluginFormatZip:
		return cp.processZip(fs, path, targetDir)
	case TypePluginFormatBin:
		return cp.processBin(fs, pluginName, path, targetDir)
	default:
		return errs.NewUserf("Unknown plugin format %s", cp.Format)
	}
}

func (cp *CustomPlugin) processBin(fs afero.Fs, name string, downloadPath string, targetDir string) error {
	target, err := Template(targetDir, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return errs.WrapUser(err, "unable to template url")
	}

	err = fs.MkdirAll(target, 0755)
	if err != nil {
		return errs.WrapUserf(err, "Could not create directory %s", target)
	}

	targetPath := path.Join(target, name)
	src, err := os.Open(downloadPath)
	if err != nil {
		return errs.WrapUserf(err, "Could not open downloaded file at %s", downloadPath)
	}

	dst, err := fs.Create(targetPath)
	if err != nil {
		return errs.WrapUserf(err, "Could not open target file at %s", targetPath)
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return errs.WrapUserf(err, "Could not move %s to %s", downloadPath, targetPath)
	}

	err = fs.Chmod(targetPath, os.FileMode(0755))
	return errs.WrapUserf(err, "Error making %s executable", targetPath)
}

// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func (cp *CustomPlugin) processTar(fs afero.Fs, path string, targetDir string) error {
	targetDir, err := Template(targetDir, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return errs.WrapUser(err, "unable to template url for custom plugin")
	}
	logrus.Debugf("untarring from %s to %s", path, targetDir)

	err = fs.MkdirAll(targetDir, 0755)
	if err != nil {
		return errs.WrapUserf(err, "Could not create directory %s", targetDir)
	}
	f, err := os.Open(path)
	if err != nil {
		return errs.WrapUser(err, "could not read staged custom plugin")
	}
	defer f.Close()
	gzr, err := gzip.NewReader(f)
	if err != nil {
		return errs.WrapUser(err, "could not create gzip reader")
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // no more files are found
		}
		if err != nil {
			return errs.WrapUser(err, "Error reading tar")
		}
		if header == nil {
			return errs.NewUser("Nil tar file header")
		}
		// the target location where the dir/file should be created
		splitTarget := strings.Split(
			filepath.Clean(header.Name),
			string(os.PathSeparator))
		// remove components if we can, otherwise skip this
		if len(splitTarget) <= cp.TarConfig.getStripComponents() {
			continue
		}
		target := filepath.Join(targetDir,
			filepath.Join(splitTarget[cp.TarConfig.getStripComponents():]...))

		switch header.Typeflag {
		case tar.TypeDir: // if its a dir and it doesn't exist create it
			err := fs.MkdirAll(target, 0755)
			if err != nil {
				return errs.WrapUserf(err, "tar: could not create directory %s", target)
			}
		case tar.TypeReg: // if it is a file create it, preserving the file mode
			destFile, err := fs.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return errs.WrapUserf(err, "tar: could not open destination file for %s", target)
			}
			_, err = io.Copy(destFile, tr)
			if err != nil {
				destFile.Close()
				return errs.WrapUserf(err, "tar: could not copy file contents")
			}
			// Manually take care of closing file since defer will pile them up
			destFile.Close()
		default:
			logrus.Warnf("tar: unrecognized tar.Type %d", header.Typeflag)
		}
	}
	return nil
}

// based on https://golangcode.com/create-zip-files-in-go/
func (cp *CustomPlugin) processZip(fs afero.Fs, downloadPath string, targetDir string) error {
	targetDir, err := Template(targetDir, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return errs.WrapUserf(err, "could not template targetDir")
	}

	err = fs.MkdirAll(targetDir, 0755)
	if err != nil {
		return errs.WrapUserf(err, "Could not create directory %s", targetDir)
	}

	r, err := zip.OpenReader(downloadPath)
	if err != nil {
		return errs.WrapUser(err, "could not read staged custom plugin")
	}
	defer r.Close()

	for _, f := range r.File {
		// We run this in a closure to invoke the `defer`s after each iteration
		err = func() error {
			rc, err := f.Open()
			if err != nil {
				return errs.WrapUserf(err, "error reading file from zip %s", f.Name)
			}
			defer rc.Close()
			fpath := filepath.Join(targetDir, f.Name)

			// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
			if !strings.HasPrefix(fpath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
				return errs.NewUserf("%s: illegal file path", fpath)
			}

			if f.FileInfo().IsDir() {
				err = os.MkdirAll(fpath, os.ModePerm)
				if err != nil {
					return errs.WrapUser(err, "zip: could not mkdirs")
				}
			} else {
				destFile, err := fs.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(f.Mode()))
				if err != nil {
					return errs.WrapUserf(err, "zip: could not open destination file for %s", fpath)
				}
				defer destFile.Close()
				_, err = io.Copy(destFile, rc)
				if err != nil {
					return errs.WrapUserf(err, "zip: could not copy file contents")
				}
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
