package util

import (
	"crypto/sha256"
	"os"
	"path/filepath"

	"github.com/chanzuckerberg/fogg/errs"
	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

func DownloadModule(fs afero.Fs, cacheDir, source string) (string, error) {

	// We want to do these operations from the root of our working repository.
	// In the case where we have a BaseFs we pull out its root. Otherwise use `pwd`.
	var pwd string
	if baseFs, ok := fs.(*afero.BasePathFs); ok {
		pwd = afero.FullBaseFsPath(baseFs, ".")
	} else {
		var e error
		pwd, e = os.Getwd()
		if e != nil {
			return "", errs.WrapUser(e, "could not get pwd")
		}
	}

	s, e := getter.Detect(source, pwd, getter.Detectors)
	if e != nil {
		return "", errs.WrapUser(e, "could not detect module type")
	}

	storage := &getter.FolderStorage{
		StorageDir: cacheDir,
	}
	h := sha256.New()
	_, e = h.Write([]byte(VersionCacheKey()))
	if e != nil {
		return "", errs.WrapUser(e, "could not hash")
	}
	_, e = h.Write([]byte(source))
	if e != nil {
		return "", errs.WrapUser(e, "could not hash")
	}
	hash := string(h.Sum(nil))

	e = storage.Get(hash, s, false)
	if e != nil {
		return "", errs.WrapUser(e, "unable to read module from local storage")
	}
	d, _, e := storage.Dir(hash)
	if e != nil {
		return "", errs.WrapUser(e, "could not get module storage dir")
	}
	return d, nil
}

func DownloadAndParseModule(fs afero.Fs, mod string) (*config.Config, error) {
	homedir, e := homedir.Dir()
	if e != nil {
		return nil, errs.WrapUser(e, "unable to find homedir")
	}

	dir := filepath.Join(homedir, ".fogg", "cache")

	d, e := DownloadModule(fs, dir, mod)
	if e != nil {
		return nil, errs.WrapUser(e, "unable to download module")
	}
	return config.LoadDir(d)
}
