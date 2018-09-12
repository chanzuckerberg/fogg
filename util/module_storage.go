package util

import (
	"crypto/sha256"
	"os"
	"path/filepath"

	"github.com/chanzuckerberg/fogg/errs"
	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform/config"
	homedir "github.com/mitchellh/go-homedir"
)

func DownloadModule(cacheDir, source string) (string, error) {
	pwd, e := os.Getwd()
	if e != nil {
		return "", errs.WrapUser(e, "could not get pwd")
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

func DownloadAndParseModule(mod string) (*config.Config, error) {
	homedir, e := homedir.Dir()
	if e != nil {
		return nil, errs.WrapUser(e, "unable to find homedir")
	}

	dir := filepath.Join(homedir, ".fogg", "cache")

	d, e := DownloadModule(dir, mod)
	if e != nil {
		return nil, errs.WrapUser(e, "unable to download module")
	}
	return config.LoadDir(d)
}
