package apply

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRemoveExtension(t *testing.T) {
	x := removeExtension("foo")
	assert.Equal(t, "foo", x)

	x = removeExtension("foo.txt")
	assert.Equal(t, "foo", x)

	x = removeExtension("foo.txt.asdf")
	assert.Equal(t, "foo.txt", x)
}

func TestApplyTemplateBasic(t *testing.T) {
	sourceFile := strings.NewReader("foo")
	dest := afero.NewMemMapFs()
	path := "bar"
	overrides := struct{ Foo string }{"foo"}

	e := applyTemplate(sourceFile, dest, path, overrides)
	assert.Nil(t, e)
	f, e := dest.Open("bar")
	assert.Nil(t, e)
	r, e := ioutil.ReadAll(f)
	assert.Nil(t, e)
	assert.Equal(t, "foo", string(r))
}

func TestApplyTemplate(t *testing.T) {
	sourceFile := strings.NewReader("Hello {{.Name}}")
	dest := afero.NewMemMapFs()
	path := "hello"
	overrides := struct{ Name string }{"World"}

	e := applyTemplate(sourceFile, dest, path, overrides)
	assert.Nil(t, e)
	f, e := dest.Open("hello")
	assert.Nil(t, e)
	r, e := ioutil.ReadAll(f)
	assert.Nil(t, e)
	assert.Equal(t, "Hello World", string(r))
}

func TestTouchFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	e := touchFile(fs, "foo")
	assert.Nil(t, e)
	r, e := readFile(fs, "foo")
	assert.Nil(t, e)
	assert.Equal(t, "", r)

	fs = afero.NewMemMapFs()

	writeFile(fs, "asdf", "jkl")

	r, e = readFile(fs, "asdf")
	assert.Nil(t, e)
	assert.Equal(t, "jkl", r)

	e = touchFile(fs, "asdf")
	assert.Nil(t, e)
	r, e = readFile(fs, "asdf")
	assert.Equal(t, "jkl", r)

}

func readFile(fs afero.Fs, path string) (string, error) {
	f, e := fs.Open(path)
	if e != nil {
		return "", e
	}
	r, e := ioutil.ReadAll(f)
	if e != nil {
		return "", e
	}
	return string(r), nil
}

func writeFile(fs afero.Fs, path string, contents string) error {
	f, e := fs.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if e != nil {
		return e
	}
	_, e = f.WriteString(contents)
	return e
}
