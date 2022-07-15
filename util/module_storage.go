package util

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/fogg/errs"
	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/kelseyhightower/envconfig"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
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

type ModuleDownloader interface {
	DownloadAndParseModule(fs afero.Fs) (*tfconfig.Module, error)
}

type Downloader struct {
	Source string
}

// Changes a URL to use the git protocol over HTTPS instead of SSH.
// If the URL was not a remote URL or an git/SSH protocol,
// it will return the path that was passed in. If it does
// convert it properly, it will add Github credentials to the path
func convertSSHToGithubHTTPURL(sURL, token string) string {
	// only detect the remote destinations
	s, err := getter.Detect(sURL, token, []getter.Detector{
		&getter.GitLabDetector{},
		&getter.GitHubDetector{},
		&getter.GitDetector{},
		&getter.BitBucketDetector{},
		&getter.S3Detector{},
		&getter.GCSDetector{},
	})
	if err != nil {
		logrus.Debug(err)
		return sURL
	}

	splits := strings.Split(s, "git::ssh://git@")
	if len(splits) != 2 {
		return sURL
	}
	u, err := url.Parse(fmt.Sprintf("https://%s", splits[1]))
	if err != nil {
		logrus.Debug(err)
		return sURL
	}
	u.User = url.User(token)

	// we want to force the git protocol
	return fmt.Sprintf("git::%s", u.String())
}

// Changes a URL to use the git protocol over HTTPS instead of SSH.
// If the URL was not a remote URL or an git/SSH protocol,
// it will return the path that was passed in. If it does
// convert it properly, it will add Github credentials to the path
func convertSSHToGithubAppHTTPURL(sURL, token string) string {
	// only detect the remote destinations
	s, err := getter.Detect(sURL, token, []getter.Detector{
		&getter.GitLabDetector{},
		&getter.GitHubDetector{},
		&getter.GitDetector{},
		&getter.BitBucketDetector{},
		&getter.S3Detector{},
		&getter.GCSDetector{},
	})
	if err != nil {
		logrus.Debug(err)
		return sURL
	}

	splits := strings.Split(s, "git::ssh://git@")
	if len(splits) != 2 {
		return sURL
	}
	u, err := url.Parse(fmt.Sprintf("https://%s", splits[1]))
	if err != nil {
		logrus.Debug(err)
		return sURL
	}
	u.User = url.UserPassword("x-access-token", token)

	// we want to force the git protocol
	return fmt.Sprintf("git::%s", u.String())
}

func MakeDownloader(src string) (*Downloader, error) {
	type HTTPAuth struct {
		GithubToken    *string
		GithubAppToken *string
	}
	var httpAuth HTTPAuth
	err := envconfig.Process("fogg", &httpAuth)
	if err != nil {
		return nil, errs.WrapUser(err, "unable to parse Github tokens")
	}
	if httpAuth.GithubToken != nil {
		src = convertSSHToGithubHTTPURL(src, *httpAuth.GithubToken)
	} else if httpAuth.GithubAppToken != nil {
		src = convertSSHToGithubAppHTTPURL(src, *httpAuth.GithubAppToken)
	}
	return &Downloader{Source: src}, nil
}

func (dd *Downloader) DownloadAndParseModule(fs afero.Fs) (*tfconfig.Module, error) {
	dir, err := GetFoggCachePath()
	if err != nil {
		return nil, err
	}
	d, err := DownloadModule(fs, dir, dd.Source)
	if err != nil {
		return nil, errs.WrapUser(err, "unable to download module")
	}
	module, diag := tfconfig.LoadModule(d)
	if diag.HasErrors() {
		return nil, errs.WrapInternal(diag.Err(), "There was an issue loading the module")
	}
	return module, nil
}
