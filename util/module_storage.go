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
	goversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/hashicorp/terraform/registry"
	"github.com/hashicorp/terraform/registry/regsrc"
	"github.com/hashicorp/terraform/registry/response"
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

func IsRegistrySourceAddr(addr string) bool {
	_, err := regsrc.ParseModuleSource(addr)
	return err == nil
}

func DownloadModule(fs afero.Fs, cacheDir, source, version string, reg *registry.Client) (string, error) {
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

	logrus.Debugf("Downloading module %q - version %q", source, version)
	if version != "" && IsRegistrySourceAddr(source) {
		logrus.Debugf("Attempting to download module %q from registry", source)
		resolvedSource, regErr := ResolveRegistryModule(source, version, reg)
		if regErr != nil {
			return "", regErr
		}
		source = resolvedSource
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

// Resolve module source based on version string
func ResolveRegistryModule(source string, version string, reg *registry.Client) (string, error) {
	var err error
	var vc goversion.Constraints
	if strings.TrimSpace(version) != "" {
		var err error
		vc, err = goversion.NewConstraint(version)
		if err != nil {
			return "", errs.WrapUser(err, fmt.Sprintf("module %q has invalid version constraint %q", source, version))
		}
	}
	// ParseModuleSource should not error because entry to this function is guarded
	addr, _ := regsrc.ParseModuleSource(source)
	hostname, _ := addr.SvcHost()
	var resp *response.ModuleVersions
	resp, err = reg.ModuleVersions(addr)
	if err != nil {
		if registry.IsModuleNotFound(err) {
			return "", errs.WrapUser(err, fmt.Sprintf("module %q cannot be found in the module registry at %s", source, hostname))
		} else {
			return "", errs.WrapUser(err, fmt.Sprintf("failed to retrieve available versions for module %q from %s", source, hostname))
		}
	}
	if len(resp.Modules) < 1 {
		return "", fmt.Errorf("the registry at %s returned an invalid response when Terraform requested available versions for module %q", hostname, source)
	}

	// The response might contain information about dependencies to potentially
	// optimize future requests, which doesn't apply here and we just take
	// the first item which is guaranteed to be the address we requested.
	modMeta := resp.Modules[0]

	var latestMatch *goversion.Version
	var latestVersion *goversion.Version
	for _, mv := range modMeta.Versions {
		v, err := goversion.NewVersion(mv.Version)
		if err != nil {
			logrus.Infof("The registry at %s returned an invalid version string %q for module %q, which Terraform ignored.", hostname, mv.Version, source)
			continue
		}

		// If we've found a pre-release version then we'll ignore it unless
		// it was exactly requested.
		if v.Prerelease() != "" && vc.String() != v.String() {
			logrus.Infof("ignoring %s for module %q because it is a pre-release and was not requested exactly", v, source)
			continue
		}

		if latestVersion == nil || v.GreaterThan(latestVersion) {
			latestVersion = v
		}

		if vc.Check(v) {
			if latestMatch == nil || v.GreaterThan(latestMatch) {
				latestMatch = v
			}
		}
	}

	if latestVersion == nil {
		return "", fmt.Errorf("module %q has no versions available on %s", addr, hostname)
	}

	if latestMatch == nil {
		return "", fmt.Errorf("there is no available version of module %q which matches the given version constraint %q. The newest available version is %s", addr, version, latestVersion)
	}

	dlAddr, err := reg.ModuleLocation(addr, latestMatch.String())
	if err != nil {
		return "", errs.WrapUser(err, fmt.Sprintf("failed to retrieve a download URL for %s %s from %s", addr, latestMatch, hostname))
	}
	source, _ = getter.SourceDirSubdir(dlAddr)
	return source, nil
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
	// Source to download Module from
	Source string
	// Version constraint string, if empty this will be ignored
	Version string
	// Terraform Registry Client to look up module based on version string
	RegistryClient *registry.Client
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
// This function should be used when using Github App credentials
// since it requires the use of "x-access-token", as the user
// which is not required by a Github Personal Access token.
// This function returns a string to use and
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
	return fmt.Sprintf("git::%s", u)
}

func MakeDownloader(src, version string, reg *registry.Client) (*Downloader, error) {
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
	return &Downloader{Source: src, Version: version, RegistryClient: reg}, nil
}

func (dd *Downloader) DownloadAndParseModule(fs afero.Fs) (*tfconfig.Module, error) {
	dir, err := GetFoggCachePath()
	if err != nil {
		return nil, err
	}
	d, err := DownloadModule(fs, dir, dd.Source, dd.Version, dd.RegistryClient)
	if err != nil {
		return nil, errs.WrapUser(err, "unable to download module")
	}
	module, diag := tfconfig.LoadModule(d)
	if diag.HasErrors() {
		return nil, errs.WrapInternal(diag.Err(), "There was an issue loading the module")
	}
	return module, nil
}
