package plugins

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/peterbourgon/diskv"
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
	URL       string            `validate:"required"`
	Format    TypePluginFormat  `validate:"required"`
	TarConfig *TarConfig        `yaml:"tar_config,omitempty"`
	TargetDir string            `yaml:"target_dir,omitempty"`
	OsMap     map[string]string `yaml:"os_map,omitempty"`
	ArchMap   map[string]string `yaml:"arch_map,omitempty"`

	cache *diskv.Diskv
}

// TarConfig configures the tar unpacking
type TarConfig struct {
	StripComponents int `yaml:"strip_components,omitempty"`
}

func (tc *TarConfig) getStripComponents() int {
	if tc == nil {
		return 0
	}
	return tc.StripComponents
}

// Install installs the custom plugin
func (cp *CustomPlugin) Install(fs afero.Fs, pluginName string) error {
	if cp.cache == nil {
		return errs.NewInternal("download cache not configured")
	}
	osName := runtime.GOOS
	if cp.OsMap != nil {
		if mappedOsName, ok := cp.OsMap[osName]; ok {
			osName = mappedOsName
		} else {
			logrus.Warnf("os %s not found in os_map", osName)
		}
	}
	archName := runtime.GOARCH
	if cp.ArchMap != nil {
		if mappedArchName, ok := cp.ArchMap[archName]; ok {
			archName = mappedArchName
		} else {
			logrus.Warnf("arch %s not found in arch_map", archName)
		}
	}
	return cp.install(fs, pluginName, cp.URL, osName, archName)
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

	targetDir, err := Template(cp.TargetDir, pluginOS, pluginArch)
	if err != nil {
		return err
	}

	fullURL, err := Template(url, pluginOS, pluginArch)
	if err != nil {
		return err
	}

	file, err := cp.fetch(fullURL)
	if err != nil {
		return err
	}
	defer file.Close()
	return cp.process(fs, pluginName, file, targetDir)
}

// Template templatizes a url with os and arch information
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

// WithTargetPath sets the target path for this plugin
func (cp *CustomPlugin) WithTargetPath(path string) *CustomPlugin {
	cp.TargetDir = path
	return cp
}

// WithCache adds a cache to plugins
func (cp *CustomPlugin) WithCache(cache *diskv.Diskv) *CustomPlugin {
	cp.cache = cache
	return cp
}

// fetch fetches the custom plugin at URL
func (cp *CustomPlugin) fetch(url string) (io.ReadCloser, error) {
	h := sha256.New()
	_, err := h.Write([]byte(url))
	if err != nil {
		return nil, errs.WrapUser(err, "error hashing url")
	}
	cacheKey := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if !cp.cache.Has(cacheKey) {
		logrus.Debugf("downloading %s to cache", url)
		resp, err := http.Get(url)
		if err != nil {
			return nil, errs.WrapUserf(err, "could not get %s", url) // FIXME
		}
		defer resp.Body.Close()
		err = cp.cache.WriteStream(cacheKey, resp.Body, true)
		if err != nil {
			return nil, errs.WrapUser(err, "could not write file to cache")
		}
	}

	f, err := cp.cache.ReadStream(cacheKey, false)
	return f, errs.WrapUserf(err, "could not read %s from cache", url)
}

// process the custom plugin
func (cp *CustomPlugin) process(fs afero.Fs, pluginName string, file io.Reader, targetDir string) error {
	switch cp.Format {
	case TypePluginFormatTar:
		return cp.processTar(fs, file, targetDir)
	case TypePluginFormatZip:
		return cp.processZip(fs, file, targetDir)
	case TypePluginFormatBin:
		return cp.processBin(fs, pluginName, file, targetDir)
	default:
		return errs.NewUserf("Unknown plugin format %s", cp.Format)
	}
}

func (cp *CustomPlugin) processBin(fs afero.Fs, name string, file io.Reader, targetDir string) error {
	err := fs.MkdirAll(targetDir, 0755)
	if err != nil {
		return errs.WrapUserf(err, "Could not create directory %s", targetDir)
	}

	targetPath := path.Join(targetDir, name)
	dst, err := fs.Create(targetPath)
	if err != nil {
		return errs.WrapUserf(err, "Could not open target file at %s", targetPath)
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		return errs.WrapUserf(err, "Could not move plugin to %s", targetPath)
	}

	err = fs.Chmod(targetPath, os.FileMode(0755))
	return errs.WrapUserf(err, "Error making %s executable", targetPath)
}

// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func (cp *CustomPlugin) processTar(fs afero.Fs, reader io.Reader, targetDir string) error {
	err := fs.MkdirAll(targetDir, 0755)
	if err != nil {
		return errs.WrapUserf(err, "Could not create directory %s", targetDir)
	}
	gzr, err := gzip.NewReader(reader)
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
			logrus.Debugf("processing tar - dir: %s", target)
			err := fs.MkdirAll(target, 0755)
			if err != nil {
				return errs.WrapUserf(err, "tar: could not create directory %s", target)
			}
		case tar.TypeReg: // if it is a file create it, preserving the file mode
			logrus.Debugf("processing tar - file: %s", target)
			dir := filepath.Dir(target)
			err := fs.MkdirAll(dir, 0755)
			if err != nil {
				return errs.WrapUserf(err, "tar: could not create directory %s", target)
			}
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
func (cp *CustomPlugin) processZip(fs afero.Fs, reader io.Reader, targetDir string) error {
	buff := bytes.NewBuffer(nil)
	size, err := io.Copy(buff, reader)
	if err != nil {
		return errs.WrapUser(err, "could not read plugin contents")
	}
	zipBytes := bytes.NewReader(buff.Bytes())
	err = fs.MkdirAll(targetDir, 0755)
	if err != nil {
		return errs.WrapUserf(err, "Could not create directory %s", targetDir)
	}

	r, err := zip.NewReader(zipBytes, size)
	if err != nil {
		return errs.WrapUser(err, "could not read staged custom plugin")
	}

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
				destFile, err := fs.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, f.Mode())
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

// GetPluginCache returns the cache used for plugins
func GetPluginCache(cacheDir string) *diskv.Diskv {
	return diskv.New(diskv.Options{
		BasePath:    cacheDir,
		Transform:   func(k string) []string { return []string{"custom_plugins", util.VersionCacheKey(), k} },
		Compression: diskv.NewGzipCompression(),
	})
}
