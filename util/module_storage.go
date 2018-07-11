package util

import (
	"crypto/sha256"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

func DownloadModule(cacheDir, source string) (string, error) {
	pwd, _ := os.Getwd()
	s, _ := getter.Detect(source, pwd, getter.Detectors)

	storage := &getter.FolderStorage{
		StorageDir: cacheDir,
	}
	h := sha256.New()
	h.Write([]byte(VersionCacheKey()))
	h.Write([]byte(source))
	hash := string(h.Sum(nil))

	storage.Get(hash, s, false)
	d, _, _ := storage.Dir(hash)
	return d, nil
}

func DownloadAndParseModule(mod string) (*config.Config, error) {
	homedir, e := homedir.Dir()
	if e != nil {
		return nil, errors.Wrap(e, "unable to find homedir")
	}

	dir := filepath.Join(homedir, ".fogg", "cache")

	d, e := DownloadModule(dir, mod)
	if e != nil {
		return nil, errors.Wrap(e, "unable to download module")
	}
	return config.LoadDir(d)
}
