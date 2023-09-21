//go:build !offline
// +build !offline

package util

import (
	"fmt"
	"os"
	"testing"

	getter "github.com/hashicorp/go-getter"
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

func TestMakeDownloader(t *testing.T) {
	r := require.New(t)
	creds := "REDACTED"
	downloader, err := MakeDownloader("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0")
	r.NoError(err)
	r.Equal("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", downloader.Source)

	os.Setenv("FOGG_GITHUBTOKEN", creds)
	defer os.Unsetenv("FOGG_GITHUBTOKEN")
	downloader, err = MakeDownloader("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0")
	r.NoError(err)
	r.Equal(fmt.Sprintf("git::https://%s@github.com/chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", creds), downloader.Source)
}

func TestMakeDownloaderGithubApp(t *testing.T) {
	r := require.New(t)
	creds := "REDACTED"
	downloader, err := MakeDownloader("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0")
	r.NoError(err)
	r.Equal("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", downloader.Source)

	os.Setenv("FOGG_GITHUBAPPTOKEN", creds)
	defer os.Unsetenv("FOGG_GITHUBAPPTOKEN")
	downloader, err = MakeDownloader("git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0")
	r.NoError(err)
	r.Equal(fmt.Sprintf("git::https://x-access-token:%s@github.com/chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", creds), downloader.Source)
}

func TestRedactURL(t *testing.T) {
	r := require.New(t)

	rurl := redactCredentials("git::https://x-access-token:1234@github.com/chanzuckerberg/shared-infra")
	r.Equal("git::https://REDACTED:REDACTED@github.com/chanzuckerberg/shared-infra", rurl)
}

func TestConvertSSHToHTTP(t *testing.T) {
	r := require.New(t)

	type test struct {
		in  string
		out string
	}
	creds := "REDACTED"
	tests := func(token string) []test {
		return []test{
			{
				in:  "git@github.com:chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0",
				out: fmt.Sprintf("git::https://%s@github.com/chanzuckerberg/test-repo//terraform/modules/eks-airflow?ref=v0.80.0", token),
			},
		}
	}(creds)
	for _, test := range tests {
		u := convertSSHToGithubHTTPURL(test.in, creds)
		s, err := getter.Detect(u, creds, []getter.Detector{
			&getter.GitHubDetector{},
		})
		r.NoError(err)
		r.Equal(test.out, u)
		r.Equal(test.out, s)
	}
}
func TestConvertSSHToHTTPFail(t *testing.T) {
	r := require.New(t)

	type test struct {
		in string
	}

	tests := func() []test {
		return []test{
			{
				in: "github.com/scholzj/terraform-aws-vpc",
			},
			{
				in: "terraform/modules/test",
			},
			{
				in: "github.com/hashicorp/go-getter?ref=abcd12",
			},
		}
	}()
	for _, test := range tests {
		s := convertSSHToGithubHTTPURL(test.in, "")
		// make sure nothing changed in a failure case
		r.Equal(test.in, s)
	}
}
