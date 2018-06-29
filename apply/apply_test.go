package apply

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/chanzuckerberg/fogg/templates"
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
	assert.Nil(t, e)
	assert.Equal(t, "jkl", r)

}

func TestCreateFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	// create new file

	e := createFile(fs, "foo", strings.NewReader("bar"))
	assert.Nil(t, e)

	r, e := readFile(fs, "foo")
	assert.Nil(t, e)
	assert.Equal(t, "bar", r)

	// not overwrite existing file

	fs = afero.NewMemMapFs()

	e = createFile(fs, "foo", strings.NewReader("bar"))
	assert.Nil(t, e)

	r, e = readFile(fs, "foo")
	assert.Nil(t, e)
	assert.Equal(t, "bar", r)

	e = createFile(fs, "foo", strings.NewReader("BAM"))
	assert.Nil(t, e)

	r, e = readFile(fs, "foo")
	assert.Nil(t, e)
	assert.Equal(t, "bar", r)

}

func TestApplySmokeTest(t *testing.T) {
	fs := afero.NewMemMapFs()
	config := `
{
  "defaults": {
    "aws_region": "reg",
    "aws_profile": "prof",
    "infra_s3_bucket": "buck",
    "project": "proj",
    "terraform_version": "0.100.0",
    "owner": "foo@example.com"
  },
  "accounts": {
    "foo": {
      "account_id": 123
    },
    "bar": {
      "account_id": 456
    }
  },
  "modules": {
    "my_module": {}
  },
  "envs": {
    "staging":{
        "components": {
            "comp1": {},
            "comp2": {}
        }
    },
    "prod": {}
  }
}
`
	afero.WriteReader(fs, "fogg.json", strings.NewReader(config))

	e := Apply(fs, "fogg.json", templates.Templates)
	assert.Nil(t, e)
}

func TestDownloadModule(t *testing.T) {
	dir, e := ioutil.TempDir("", "fogg")
	assert.Nil(t, e)

	e = downloadModule(dir, "./test-module")
	assert.Nil(t, e)
	// TODO more asserts
}

func TestDownloadAndParseModule(t *testing.T) {
	c, e := downloadAndParseModule("./test-module")
	assert.Nil(t, e)
	assert.NotNil(t, c)
	assert.Len(t, c.Variables, 2)
	assert.Len(t, c.Outputs, 2)
}

func TestApplyModule(t *testing.T) {
	fs := afero.NewMemMapFs()

	e := applyModule(fs, "mymodule", "./test-module", templates.Templates.Module)
	assert.Nil(t, e)

	s, e := fs.Stat("mymodule")
	assert.Nil(t, e)
	assert.True(t, s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	assert.Nil(t, e)
	r, e := afero.ReadFile(fs, "mymodule/main.tf")
	assert.Nil(t, e)
	expected := `module "test-module" {
  source = "./test-module"
  foo = "${var.foo}"
  bar = "${var.bar}"
  
}`
	assert.Equal(t, expected, string(r))

	_, e = fs.Stat("mymodule/outputs.tf")
	assert.Nil(t, e)
	r, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	assert.Nil(t, e)
	expected = `output "foo" {
  value = "${module.test-module.foo}"
}

output "bar" {
  value = "${module.test-module.bar}"
}

`
	assert.Equal(t, expected, string(r))
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
