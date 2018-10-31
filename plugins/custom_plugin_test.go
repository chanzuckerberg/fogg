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
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCustomPluginTar(t *testing.T) {
	a := assert.New(t)
	pluginName := "test-provider"
	fs := afero.NewMemMapFs()

	files := []string{"test.txt", "terraform-provider-testing"}
	tarPath := generateTar(t, files)
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
	customPlugin.SetTargetPath(plugins.BinDir)
	a.Nil(customPlugin.Install(fs, pluginName))

	afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		a.Nil(err)
		return nil
	})

	for _, file := range files {
		filePath := path.Join(plugins.BinDir, file)
		fi, err := fs.Stat(filePath)
		a.Nil(err)
		a.False(fi.IsDir())
		a.Equal(fi.Mode(), os.FileMode(0664))

		bytes, err := afero.ReadFile(fs, filePath)
		a.Nil(err)
		a.Equal(bytes, []byte(file)) // We wrote the filename as the contents as well
	}
}

func TestCustomPluginZip(t *testing.T) {
	a := assert.New(t)
	pluginName := "test-provider"
	fs := afero.NewMemMapFs()

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
	customPlugin.SetTargetPath(plugins.BinDir)
	a.Nil(customPlugin.Install(fs, pluginName))

	afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		a.Nil(err)
		return nil
	})

	for _, file := range files {
		filePath := path.Join(plugins.BinDir, file)
		fi, err := fs.Stat(filePath)
		a.Nil(err)
		a.False(fi.IsDir())
		a.Equal(fi.Mode(), os.FileMode(0664))

		bytes, err := afero.ReadFile(fs, filePath)
		a.Nil(err)
		a.Equal(bytes, []byte(file)) // We wrote the filename as the contents as well
	}
}

func TestCustomPluginBin(t *testing.T) {
	a := assert.New(t)
	pluginName := "test-provider"
	fs := afero.NewMemMapFs()
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

	customPlugin.SetTargetPath(plugins.BinDir)
	a.Nil(customPlugin.Install(fs, pluginName))

	afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		a.Nil(err)
		return nil
	})

	customPluginPath := path.Join(plugins.BinDir, pluginName)
	f, err := fs.Open(customPluginPath)
	a.Nil(err)

	contents, err := ioutil.ReadAll(f)
	a.Nil(err)
	a.Equal(string(contents), fileContents)

	fi, err := fs.Stat(customPluginPath)
	a.Nil(err)
	a.False(fi.IsDir())
	a.Equal(fi.Mode().Perm(), os.FileMode(0755))
}

func generateTar(t *testing.T, files []string) string {
	a := assert.New(t)

	f, err := ioutil.TempFile("", "testing")
	a.Nil(err)
	defer f.Close()
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		header := new(tar.Header)
		header.Name = file
		header.Size = int64(len([]byte(file)))
		header.Mode = int64(0664)
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
		header.SetMode(os.FileMode(0664))
		writer, err := zipWriter.CreateHeader(header)
		a.Nil(err)
		_, err = io.Copy(writer, strings.NewReader(file))
		a.Nil(err)
	}
	return f.Name()
}
