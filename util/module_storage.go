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

func redactCredentials(source string) string {
	splits := strings.Split(source, "git::")
	if len(splits) != 2 {
		return source
	}
	u, err := url.Parse(splits[1])
	if err != nil {
		return source
	}
	u.User = url.UserPassword("REDACTED", "REDACTED")
	return fmt.Sprintf("git::%s", u)
}

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
	// Sometimes, like in a Github App credential case,
	// the source will have dynamic 1-time credentials
	// Make sure to remove these so we don't break the cache.
	_, err = h.Write([]byte(redactCredentials(source)))
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

var _gitDetectors = []getter.Detector{
	&getter.GitLabDetector{},
	&getter.GitHubDetector{},
	&getter.GitDetector{},
	&getter.BitBucketDetector{},
	&getter.S3Detector{},
	&getter.GCSDetector{},
}

// injectGitHubCredentials normalizes a module source URL and injects GitHub
// credentials into it. Handles SSH (git@github.com:…), HTTPS
// (git::https://github.com/…), and GitHub shorthand (github.com/…) formats.
// Non-GitHub and local-path sources are returned unchanged.
func injectGitHubCredentials(sURL string, userInfo *url.Userinfo) string {
	s, err := getter.Detect(sURL, "", _gitDetectors)
	if err != nil {
		logrus.Debug(err)
		return sURL
	}

	if !strings.Contains(s, "github.com") {
		return sURL
	}

	// SSH: git::ssh://git@github.com/owner/repo…
	if parts := strings.SplitN(s, "git::ssh://git@", 2); len(parts) == 2 {
		u, err := url.Parse("https://" + parts[1])
		if err != nil {
			logrus.Debug(err)
			return sURL
		}
		u.User = userInfo
		return fmt.Sprintf("git::%s", u.String())
	}

	// HTTPS: git::https://github.com/owner/repo… (with or without existing creds)
	if parts := strings.SplitN(s, "git::", 2); len(parts) == 2 {
		u, err := url.Parse(parts[1])
		if err != nil {
			logrus.Debug(err)
			return sURL
		}
		if u.Scheme == "https" {
			u.User = userInfo
			return fmt.Sprintf("git::%s", u.String())
		}
	}

	return sURL
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
		src = injectGitHubCredentials(src, url.User(*httpAuth.GithubToken))
	} else if httpAuth.GithubAppToken != nil {
		src = injectGitHubCredentials(src, url.UserPassword("x-access-token", *httpAuth.GithubAppToken))
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
