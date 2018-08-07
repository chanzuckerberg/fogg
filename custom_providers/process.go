package customproviders

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

// Install installs the custom provider
func (cp *CustomProvider) Install(providerName string, dest afero.Fs) error {
	if cp == nil {
		return errors.New("nil CustomProvider")
	}
	if dest == nil {
		return errors.New("nil fs")
	}

	tmpPath, err := cp.fetch(providerName)
	if err != nil {
		return err
	}
	return cp.process(tmpPath, dest)
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
func (cp *CustomProvider) process(path string, dest afero.Fs) error {
	switch cp.Format {
	// Binary formats do not need any special processing
	case TypeProviderFormatBinary:
		log.Debug("Binary format, nothing to do")
		return nil
	case TypeProviderFormatTar:
		return cp.processTar(path, dest)
	default:
		return errors.Errorf("Unknown provider format %s", cp.Format)
	}
}

// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func (cp *CustomProvider) processTar(path string, dest afero.Fs) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "custom provider")
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
		target := filepath.Join(".terraform.d/plugin-cache", header.Name) // TODO: don't hard code this
		switch header.Typeflag {
		case tar.TypeDir: // if its a dir and it doesn't exist create it
			if _, err := dest.Stat(target); err != nil {
				if err := dest.MkdirAll(target, 0755); err != nil {
					return errors.Wrapf(err, "tar: could not create directory for %s", target)
				}
			}
		case tar.TypeReg: // if it is a file create it
			destFile, err := dest.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
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
		}
	}

	return nil
}
