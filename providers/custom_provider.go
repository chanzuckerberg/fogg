package providers

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// TypeProviderFormat is the provider format such as binary, zip, tar
type TypeProviderFormat string

const (
	// TypeProviderFormatTar is a tar archived provider
	TypeProviderFormatTar TypeProviderFormat = "tar"
)

// CustomProvider is a custom terraform provider
type CustomProvider struct {
	URL    string             `json:"url" validate:"required"`
	Format TypeProviderFormat `json:"format" validate:"required"`
}

// Install installs the custom provider
func (cp *CustomProvider) Install(fs afero.Fs, providerName string) error {
	if cp == nil {
		return errors.New("nil CustomProvider")
	}
	if fs == nil {
		return errors.New("nil fs")
	}

	tmpPath, err := cp.fetch(providerName)
	defer os.Remove(tmpPath) // clean up
	if err != nil {
		return err
	}
	return cp.process(fs, providerName, tmpPath)
}

// fetch fetches the custom provider at URL
func (cp *CustomProvider) fetch(providerName string) (string, error) {
	tmpFile, err := ioutil.TempFile("", providerName)
	if err != nil {
		return "", errors.Wrap(err, "could not create temporary directory")
	}
	resp, err := http.Get(cp.URL)
	if err != nil {
		return "", errors.Wrapf(err, "could not get %s", cp.URL)
	}
	defer resp.Body.Close()
	_, err = io.Copy(tmpFile, resp.Body)
	return tmpFile.Name(), errors.Wrap(err, "could not download file")
}

// process the custom provider
func (cp *CustomProvider) process(fs afero.Fs, providerName string, path string) error {
	switch cp.Format {
	case TypeProviderFormatTar:
		return cp.processTar(fs, path)
	default:
		return errors.Errorf("Unknown provider format %s", cp.Format)
	}
}

// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func (cp *CustomProvider) processTar(fs afero.Fs, path string) error {
	err := fs.MkdirAll(CustomPluginCacheDir, 0755)
	if err != nil {
		return errors.Wrapf(err, "Could not create directory %s", CustomPluginCacheDir)
	}
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "could not read staged custom provider")
	}
	defer f.Close()
	gzr, err := gzip.NewReader(f)
	if err != nil {
		return errors.Wrap(err, "could not create gzip reader")
	}
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // no more files are found
		}
		if err != nil {
			return errors.Wrap(err, "Error reading tar")
		}
		if header == nil {
			return errors.New("Nil tar file header")
		}
		// the target location where the dir/file should be created
		target := filepath.Join(CustomPluginCacheDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir: // if its a dir and it doesn't exist create it
			err := fs.MkdirAll(target, 0755)
			if err != nil {
				return errors.Wrapf(err, "tar: could not create directory %s", target)
			}
		case tar.TypeReg: // if it is a file create it
			destFile, err := fs.Create(target)
			if err != nil {
				return errors.Wrapf(err, "tar: could not open destination file for %s", target)
			}
			_, err = io.Copy(destFile, tr)
			if err != nil {
				destFile.Close()
				return errors.Wrapf(err, "tar: could not copy file contents")
			}
			// Manually take care of closing file since defer will pile them up
			destFile.Close()
		default:
			log.Warnf("tar: unrecognized tar.Type %d", header.Typeflag)
		}
	}
	return nil
}
