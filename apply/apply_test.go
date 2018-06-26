package apply

import (
	"io/ioutil"
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
