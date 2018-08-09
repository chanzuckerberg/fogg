package providers_test

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/chanzuckerberg/fogg/providers"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCustomProviderTar(t *testing.T) {
	a := assert.New(t)
	providerName := "test-provider"
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

	customProvider := &providers.CustomProvider{
		URL:    ts.URL,
		Format: providers.TypeProviderFormatTar,
	}

	err := customProvider.Install(providerName, fs)
	a.Nil(err)

	afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		a.Nil(err)
		return nil
	})

	for _, file := range files {
		filePath := path.Join(providers.CustomPluginCacheDir, file)
		fi, err := fs.Stat(filePath)
		a.Nil(err)
		a.False(fi.IsDir())

		bytes, err := afero.ReadFile(fs, filePath)
		a.Nil(err)
		a.Equal(bytes, []byte(file)) // We wrote the filename as the contents as well
	}
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
		header.Mode = int64(0644)
		header.Typeflag = tar.TypeReg

		a.Nil(tw.WriteHeader(header))
		_, err = fmt.Fprintf(tw, file)
		a.Nil(err)
	}

	return f.Name()
}
