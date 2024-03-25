package plugins_test

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"

	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCustomPluginTar(t *testing.T) {
	r := require.New(t)
	pluginName := "test-provider"
	fs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	cacheDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(cacheDir)
	cache := plugins.GetPluginCache(cacheDir)

	files := []string{"test.txt", "terraform-provider-testing"}
	tarPath := generateTar(t, files, nil)
	defer os.Remove(tarPath)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		f, err := os.Open(tarPath)
		r.Nil(err)
		_, err = io.Copy(w, f)
		r.Nil(err)
	}))
	defer ts.Close()

	customPlugin := &plugins.CustomPlugin{
		URL:    ts.URL,
		Format: plugins.TypePluginFormatTar,
	}
	customPlugin.WithTargetPath(plugins.CustomPluginDir).WithCache(cache)
	r.Nil(customPlugin.Install(fs, pluginName))

	r.NoError(afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		r.Nil(err)
		return nil
	}))

	for _, file := range files {
		filePath := path.Join(plugins.CustomPluginDir, file)
		fi, err := fs.Stat(filePath)
		r.Nil(err)
		r.False(fi.IsDir())
		r.Equal(fi.Mode(), os.FileMode(0755))

		bytes, err := afero.ReadFile(fs, filePath)
		r.Nil(err)
		r.Equal(bytes, []byte(file)) // We wrote the filename as the contents as well
	}
}

func TestCustomPluginTarStripComponents(t *testing.T) {
	r := require.New(t)
	pluginName := "test-provider"
	fs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	cacheDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(cacheDir)
	cache := plugins.GetPluginCache(cacheDir)

	files := []string{"a/test.txt", "terraform-provider-testing"}
	expectedFiles := []string{"test.txt", ""}
	dirs := []string{"a"}
	tarPath := generateTar(t, files, dirs)
	defer os.Remove(tarPath)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		f, err := os.Open(tarPath)
		r.Nil(err)
		_, err = io.Copy(w, f)
		r.Nil(err)
	}))
	defer ts.Close()

	customPlugin := &plugins.CustomPlugin{
		URL:    ts.URL,
		Format: plugins.TypePluginFormatTar,
		TarConfig: &plugins.TarConfig{
			StripComponents: 1,
		},
	}
	customPlugin.WithTargetPath(plugins.CustomPluginDir).WithCache(cache)
	r.Nil(customPlugin.Install(fs, pluginName))

	r.NoError(afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		r.Nil(err)
		return nil
	}))

	for idx, file := range expectedFiles {
		// files we expect to skip
		if file == "" {
			filePath := path.Join(plugins.CustomPluginDir, files[idx])
			_, err := fs.Stat(filePath)
			r.NotNil(err)
			r.True(os.IsNotExist(err))
			continue
		}
		filePath := path.Join(plugins.CustomPluginDir, file)
		fi, err := fs.Stat(filePath)
		r.Nil(err)
		r.False(fi.IsDir())
		r.Equal(fi.Mode(), os.FileMode(0755))

		bytes, err := afero.ReadFile(fs, filePath)
		r.Nil(err)
		r.Equal(bytes, []byte(files[idx])) // We wrote the filename as the contents as well
	}
}
func TestCustomPluginZip(t *testing.T) {
	r := require.New(t)
	pluginName := "test-provider"
	fs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	cacheDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(cacheDir)
	cache := plugins.GetPluginCache(cacheDir)

	files := []string{"test.txt", "terraform-provider-testing"}
	zipPath := generateZip(t, files)
	defer os.Remove(zipPath)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		f, err := os.Open(zipPath)
		r.Nil(err)
		_, err = io.Copy(w, f)
		r.Nil(err)
	}))
	defer ts.Close()

	customPlugin := &plugins.CustomPlugin{
		URL:    ts.URL,
		Format: plugins.TypePluginFormatZip,
	}
	customPlugin.WithTargetPath(plugins.CustomPluginDir).WithCache(cache)
	r.Nil(customPlugin.Install(fs, pluginName))

	r.NoError(afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		r.Nil(err)
		return nil
	}))

	for _, file := range files {
		filePath := path.Join(plugins.CustomPluginDir, file)
		fi, err := fs.Stat(filePath)
		r.Nil(err)
		r.False(fi.IsDir())
		r.Equal(fi.Mode(), os.FileMode(0755))

		bytes, err := afero.ReadFile(fs, filePath)
		r.Nil(err)
		r.Equal(bytes, []byte(file)) // We wrote the filename as the contents as well
	}
}

func TestCustomPluginBin(t *testing.T) {
	r := require.New(t)
	pluginName := "test-provider"
	fs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)
	fileContents := "some contents"

	cacheDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(cacheDir)
	cache := plugins.GetPluginCache(cacheDir)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprint(w, fileContents)
		r.Nil(err)
	}))
	defer ts.Close()

	customPlugin := &plugins.CustomPlugin{
		URL:    ts.URL,
		Format: plugins.TypePluginFormatBin,
	}

	customPlugin.WithTargetPath(plugins.TerraformCustomPluginCacheDir).WithCache(cache)
	r.Nil(customPlugin.Install(fs, pluginName))

	r.NoError(afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		r.Nil(err)
		return nil
	}))

	// Make sure we properly template
	pluginDir, err := plugins.Template(plugins.TerraformCustomPluginCacheDir, runtime.GOOS, runtime.GOARCH)
	r.NoError(err)
	customPluginPath := path.Join(pluginDir, pluginName)
	f, err := fs.Open(customPluginPath)
	r.Nil(err)

	contents, err := io.ReadAll(f)
	r.Nil(err)
	r.Equal(string(contents), fileContents)

	fi, err := fs.Stat(customPluginPath)
	r.Nil(err)
	r.False(fi.IsDir())
	r.Equal(os.FileMode(0755), fi.Mode().Perm())
}

func generateTar(t *testing.T, files []string, dirs []string) string {
	r := require.New(t)

	f, err := os.CreateTemp("", "testing")
	r.Nil(err)
	defer f.Close()
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, dir := range dirs {
		header := new(tar.Header)
		header.Name = dir
		header.Size = int64(0)
		header.Mode = int64(0777)
		header.Typeflag = tar.TypeDir

		r.Nil(tw.WriteHeader(header))
	}

	for _, file := range files {
		header := new(tar.Header)
		header.Name = file
		header.Size = int64(len([]byte(file)))
		header.Mode = int64(0755)
		header.Typeflag = tar.TypeReg

		r.Nil(tw.WriteHeader(header))
		_, err = fmt.Fprint(tw, file)
		r.Nil(err)
	}

	return f.Name()
}

// based on https://golangcode.com/create-zip-files-in-go/
func generateZip(t *testing.T, files []string) string {
	r := require.New(t)
	f, err := os.CreateTemp("", "testing")
	r.Nil(err)

	defer f.Close()

	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	for _, file := range files {
		header := &zip.FileHeader{
			Name: file,
		}
		header.SetMode(os.FileMode(0755))
		writer, err := zipWriter.CreateHeader(header)
		r.Nil(err)
		_, err = io.Copy(writer, strings.NewReader(file))
		r.Nil(err)
	}
	return f.Name()
}

func TestTemplate(t *testing.T) {
	type args struct {
		url  string
		arch string
		os   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"no-op", args{"foo", "bar", "bam"}, "foo", false},
		{"os", args{"{{.OS}}", "bar", "bam"}, "bar", false},
		{"arch", args{"{{.Arch}}", "bar", "bam"}, "bam", false},
		{"os_arch", args{"{{.OS}}_{{.Arch}}", "bar", "bam"}, "bar_bam", false},
		{"bad template", args{"{{.asdf", "bar", "bam"}, "", true},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			got, err := plugins.Template(tt.args.url, tt.args.arch, tt.args.os)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Template() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstallErrorsOnNoCache(t *testing.T) {
	r := require.New(t)
	fs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	customPlugin := &plugins.CustomPlugin{
		URL:    "my url",
		Format: plugins.TypePluginFormatBin,
	}

	err = customPlugin.Install(fs, "my plugin name")
	r.Error(err, "download cache not configured")
}
