package util

import (
	"crypto/sha256"
	"os"
	"path/filepath"

	"github.com/chanzuckerberg/fogg/errs"
	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

func DownloadModule(fs afero.Fs, cacheDir, source string) (string, error) {

	// We want to do these operations from the root of our working repository.
	// In the case where we have a BaseFs we pull out its root. Otherwise use `pwd`.
	var pwd string
	var err error
	if baseFs, ok := fs.(*afero.BasePathFs); ok {
		pwd = afero.FullBaseFsPath(baseFs, ".")
	} else {
		pwd, err = os.Getwd()
		if err != nil {
			return "", errs.WrapUser(err, "could not get pwd")
		}
	}

	s, err := getter.Detect(source, pwd, getter.Detectors)
	if err != nil {
		return "", errs.WrapUser(err, "could not detect module type")
	}

	storage := &getter.FolderStorage{
		StorageDir: cacheDir,
	}
	h := sha256.New()
	_, err = h.Write([]byte(VersionCacheKey()))
	if err != nil {
		return "", errs.WrapUser(err, "could not hash")
	}
	_, err = h.Write([]byte(source))
	if err != nil {
		return "", errs.WrapUser(err, "could not hash")
	}
	hash := string(h.Sum(nil))

	err = storage.Get(hash, s, false)
	if err != nil {
		return "", errs.WrapUser(err, "unable to read module from local storage")
	}
	d, _, err := storage.Dir(hash)
	if err != nil {
		return "", errs.WrapUser(err, "could not get module storage dir")
	}
	return d, nil
}

func GetFoggCachePath() (string, error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return "", errs.WrapUser(err, "unable to find homedir")
	}
	dir := filepath.Join(homedir, ".fogg", "cache")
	return dir, nil
}

func DownloadAndParseModule(fs afero.Fs, mod string) (*tfconfig.Module, error) {
	dir, err := GetFoggCachePath()
	if err != nil {
		return nil, err
	}
	d, err := DownloadModule(fs, dir, mod)
	if err != nil {
		return nil, errs.WrapUser(err, "unable to download module")
	}
	module, diag := tfconfig.LoadModule(d)
	if diag.HasErrors() {
		return nil, errs.WrapInternal(diag.Err(), "There was an issue loading the module")
	}
	return module, nil
}
