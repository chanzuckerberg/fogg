//go:build !offline
// +build !offline

package util

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestDownloadModule(t *testing.T) {
	r := require.New(t)
	dir, e := os.MkdirTemp("", "fogg")
	r.Nil(e)

	pwd, e := os.Getwd()
	r.NoError(e)

	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	d, e := DownloadModule(fs, dir, "github.com/chanzuckerberg/fogg-test-module")
	r.NoError(e)
	r.NotNil(d)
	r.NotEmpty(d)
	// TODO more asserts
}

func TestDownloadAndParseModule(t *testing.T) {
	r := require.New(t)

	pwd, e := os.Getwd()
	r.NoError(e)
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	downloader, err := MakeDownloader("github.com/chanzuckerberg/fogg-test-module")
	r.NoError(err)
	c, e := downloader.DownloadAndParseModule(fs)
	r.Nil(e)
	r.NotNil(c)
	r.NotNil(c.Variables)
	r.NotNil(c.Outputs)
	r.Len(c.Variables, 2)
	r.Len(c.Outputs, 2)
}

func TestMakeDownloaderSSH(t *testing.T) {
	r := require.New(t)
	creds := "REDACTED"

	// Without token, source is unchanged
	downloader, err := MakeDownloader("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0")
	r.NoError(err)
	r.Equal("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", downloader.Source)

	// PAT: SSH → HTTPS with token as user
	t.Setenv("FOGG_GITHUBTOKEN", creds)
	downloader, err = MakeDownloader("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0")
	r.NoError(err)
	r.Equal(fmt.Sprintf("git::https://%s@github.com/chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", creds), downloader.Source)
}

func TestMakeDownloaderHTTPS(t *testing.T) {
	r := require.New(t)
	creds := "REDACTED"

	// Without token, source is unchanged
	downloader, err := MakeDownloader("git::https://github.com/chanzuckerberg/cztack//aws-vpc-env?ref=aws-vpc-env-v4")
	r.NoError(err)
	r.Equal("git::https://github.com/chanzuckerberg/cztack//aws-vpc-env?ref=aws-vpc-env-v4", downloader.Source)

	// PAT: HTTPS URL gets token injected
	t.Setenv("FOGG_GITHUBTOKEN", creds)
	downloader, err = MakeDownloader("git::https://github.com/chanzuckerberg/cztack//aws-vpc-env?ref=aws-vpc-env-v4")
	r.NoError(err)
	r.Equal(fmt.Sprintf("git::https://%s@github.com/chanzuckerberg/cztack//aws-vpc-env?ref=aws-vpc-env-v4", creds), downloader.Source)
}

func TestMakeDownloaderGithubAppSSH(t *testing.T) {
	r := require.New(t)
	creds := "REDACTED"

	// Without token, source is unchanged
	downloader, err := MakeDownloader("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0")
	r.NoError(err)
	r.Equal("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", downloader.Source)

	// App token: SSH → HTTPS with x-access-token:token
	t.Setenv("FOGG_GITHUBAPPTOKEN", creds)
	downloader, err = MakeDownloader("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0")
	r.NoError(err)
	r.Equal(fmt.Sprintf("git::https://x-access-token:%s@github.com/chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", creds), downloader.Source)
}

func TestMakeDownloaderGithubAppHTTPS(t *testing.T) {
	r := require.New(t)
	creds := "REDACTED"

	// App token: HTTPS URL gets x-access-token:token injected
	t.Setenv("FOGG_GITHUBAPPTOKEN", creds)
	downloader, err := MakeDownloader("git::https://github.com/chanzuckerberg/cztack//aws-vpc-env?ref=aws-vpc-env-v4")
	r.NoError(err)
	r.Equal(fmt.Sprintf("git::https://x-access-token:%s@github.com/chanzuckerberg/cztack//aws-vpc-env?ref=aws-vpc-env-v4", creds), downloader.Source)
}

func TestRedactURL(t *testing.T) {
	r := require.New(t)

	rurl := redactCredentials("git::https://x-access-token:1234@github.com/chanzuckerberg/shared-infra")
	r.Equal("git::https://REDACTED:REDACTED@github.com/chanzuckerberg/shared-infra", rurl)
}

func TestInjectGitHubCredentials(t *testing.T) {
	creds := "MY-TOKEN"
	patUser := url.User(creds)
	appUser := url.UserPassword("x-access-token", creds)

	tests := []struct {
		name     string
		in       string
		userInfo *url.Userinfo
		want     string
	}{
		{
			name:     "SSH with PAT",
			in:       "git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0",
			userInfo: patUser,
			want:     fmt.Sprintf("git::https://%s@github.com/chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", creds),
		},
		{
			name:     "SSH with App token",
			in:       "git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0",
			userInfo: appUser,
			want:     fmt.Sprintf("git::https://x-access-token:%s@github.com/chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", creds),
		},
		{
			name:     "HTTPS with PAT",
			in:       "git::https://github.com/chanzuckerberg/cztack//aws-vpc-env?ref=aws-vpc-env-v4",
			userInfo: patUser,
			want:     fmt.Sprintf("git::https://%s@github.com/chanzuckerberg/cztack//aws-vpc-env?ref=aws-vpc-env-v4", creds),
		},
		{
			name:     "HTTPS with App token",
			in:       "git::https://github.com/chanzuckerberg/cztack//aws-vpc-env?ref=aws-vpc-env-v4",
			userInfo: appUser,
			want:     fmt.Sprintf("git::https://x-access-token:%s@github.com/chanzuckerberg/cztack//aws-vpc-env?ref=aws-vpc-env-v4", creds),
		},
		{
			name:     "HTTPS no subpath",
			in:       "git::https://github.com/chanzuckerberg/shared-infra?ref=v1.0.0",
			userInfo: patUser,
			want:     fmt.Sprintf("git::https://%s@github.com/chanzuckerberg/shared-infra?ref=v1.0.0", creds),
		},
		{
			name:     "already-credentialed HTTPS replaces creds",
			in:       "git::https://old-token@github.com/chanzuckerberg/cztack//mod?ref=v1",
			userInfo: appUser,
			want:     fmt.Sprintf("git::https://x-access-token:%s@github.com/chanzuckerberg/cztack//mod?ref=v1", creds),
		},
		{
			name:     "local path unchanged",
			in:       "terraform/modules/test",
			userInfo: patUser,
			want:     "terraform/modules/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			got := injectGitHubCredentials(tt.in, tt.userInfo)
			r.Equal(tt.want, got)
		})
	}
}
