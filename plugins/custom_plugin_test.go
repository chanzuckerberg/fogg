package plugins_test

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCustomPluginTar(t *testing.T) {
	a := assert.New(t)
	pluginName := "test-provider"
	fs, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)

	files := []string{"test.txt", "terraform-provider-testing"}
	tarPath := generateTar(t, files, nil)
	defer os.Remove(tarPath)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(tarPath)
		a.Nil(err)
		_, err = io.Copy(w, f)
		a.Nil(err)
	}))
	defer ts.Close()

	customPlugin := &plugins.CustomPlugin{
		URL:    ts.URL,
		Format: plugins.TypePluginFormatTar,
	}
	customPlugin.SetTargetPath(plugins.CustomPluginDir)
	a.Nil(customPlugin.Install(fs, pluginName))

	a.NoError(afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		a.Nil(err)
		return nil
	}))

	for _, file := range files {
		filePath := path.Join(plugins.CustomPluginDir, file)
		fi, err := fs.Stat(filePath)
		a.Nil(err)
		a.False(fi.IsDir())
		a.Equal(fi.Mode(), os.FileMode(0755))

		bytes, err := afero.ReadFile(fs, filePath)
		a.Nil(err)
		a.Equal(bytes, []byte(file)) // We wrote the filename as the contents as well
	}
}

func TestCustomPluginTarStripComponents(t *testing.T) {
	a := assert.New(t)
	pluginName := "test-provider"
	fs, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)

	files := []string{"a/test.txt", "terraform-provider-testing"}
	expected_files := []string{"test.txt", ""}
	dirs := []string{"a"}
	tarPath := generateTar(t, files, dirs)
	defer os.Remove(tarPath)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(tarPath)
		a.Nil(err)
		_, err = io.Copy(w, f)
		a.Nil(err)
	}))
	defer ts.Close()

	customPlugin := &plugins.CustomPlugin{
		URL:    ts.URL,
		Format: plugins.TypePluginFormatTar,
		TarConfig: &plugins.TarConfig{
			StripComponents: 1,
		},
	}
	customPlugin.SetTargetPath(plugins.CustomPluginDir)
	a.Nil(customPlugin.Install(fs, pluginName))

	a.NoError(afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		a.Nil(err)
		return nil
	}))

	for idx, file := range expected_files {
		// files we expect to skip
		if file == "" {
			filePath := path.Join(plugins.CustomPluginDir, files[idx])
			_, err := fs.Stat(filePath)
			a.NotNil(err)
			a.True(os.IsNotExist(err))
			continue
		}
		filePath := path.Join(plugins.CustomPluginDir, file)
		fi, err := fs.Stat(filePath)
		a.Nil(err)
		a.False(fi.IsDir())
		a.Equal(fi.Mode(), os.FileMode(0755))

		bytes, err := afero.ReadFile(fs, filePath)
		a.Nil(err)
		a.Equal(bytes, []byte(files[idx])) // We wrote the filename as the contents as well
	}
}
func TestCustomPluginZip(t *testing.T) {
	a := assert.New(t)
	pluginName := "test-provider"
	fs, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)

	files := []string{"test.txt", "terraform-provider-testing"}
	zipPath := generateZip(t, files)
	defer os.Remove(zipPath)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(zipPath)
		a.Nil(err)
		_, err = io.Copy(w, f)
		a.Nil(err)
	}))
	defer ts.Close()

	customPlugin := &plugins.CustomPlugin{
		URL:    ts.URL,
		Format: plugins.TypePluginFormatZip,
	}
	customPlugin.SetTargetPath(plugins.CustomPluginDir)
	a.Nil(customPlugin.Install(fs, pluginName))

	a.NoError(afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		a.Nil(err)
		return nil
	}))

	for _, file := range files {
		filePath := path.Join(plugins.CustomPluginDir, file)
		fi, err := fs.Stat(filePath)
		a.Nil(err)
		a.False(fi.IsDir())
		a.Equal(fi.Mode(), os.FileMode(0755))

		bytes, err := afero.ReadFile(fs, filePath)
		a.Nil(err)
		a.Equal(bytes, []byte(file)) // We wrote the filename as the contents as well
	}
}

func TestCustomPluginBin(t *testing.T) {
	a := assert.New(t)
	pluginName := "test-provider"
	fs, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)
	fileContents := "some contents"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, fileContents)
		a.Nil(err)
	}))
	defer ts.Close()

	customPlugin := &plugins.CustomPlugin{
		URL:    ts.URL,
		Format: plugins.TypePluginFormatBin,
	}

	customPlugin.SetTargetPath(plugins.CustomPluginDir)
	a.Nil(customPlugin.Install(fs, pluginName))

	a.NoError(afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		a.Nil(err)
		return nil
	}))

	customPluginPath := path.Join(plugins.CustomPluginDir, pluginName)
	f, err := fs.Open(customPluginPath)
	a.Nil(err)

	contents, err := ioutil.ReadAll(f)
	a.Nil(err)
	a.Equal(string(contents), fileContents)

	fi, err := fs.Stat(customPluginPath)
	a.Nil(err)
	a.False(fi.IsDir())
	a.Equal(os.FileMode(0755), fi.Mode().Perm())
}

func generateTar(t *testing.T, files []string, dirs []string) string {
	a := assert.New(t)

	f, err := ioutil.TempFile("", "testing")
	a.Nil(err)
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

		a.Nil(tw.WriteHeader(header))
	}

	for _, file := range files {
		header := new(tar.Header)
		header.Name = file
		header.Size = int64(len([]byte(file)))
		header.Mode = int64(0755)
		header.Typeflag = tar.TypeReg

		a.Nil(tw.WriteHeader(header))
		_, err = fmt.Fprintf(tw, file)
		a.Nil(err)
	}

	return f.Name()
}

// based on https://golangcode.com/create-zip-files-in-go/
func generateZip(t *testing.T, files []string) string {
	a := assert.New(t)
	f, err := ioutil.TempFile("", "testing")
	a.Nil(err)

	defer f.Close()

	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	for _, file := range files {

		header := &zip.FileHeader{
			Name: file,
		}
		header.SetMode(os.FileMode(0755))
		writer, err := zipWriter.CreateHeader(header)
		a.Nil(err)
		_, err = io.Copy(writer, strings.NewReader(file))
		a.Nil(err)
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
	for _, tt := range tests {
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
