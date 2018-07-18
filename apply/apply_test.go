package apply

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/templates"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}
func TestRemoveExtension(t *testing.T) {
	x := removeExtension("foo")
	assert.Equal(t, "foo", x)

	x = removeExtension("foo.txt")
	assert.Equal(t, "foo", x)

	x = removeExtension("foo.txt.asdf")
	assert.Equal(t, "foo.txt", x)

	x = removeExtension("foo/bar.txt")
	assert.Equal(t, "foo/bar", x)
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
	json := `
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
	c, e := config.ReadConfig(ioutil.NopCloser(strings.NewReader(json)))
	assert.Nil(t, e)

	e = Apply(fs, c, templates.Templates, true)
	assert.Nil(t, e)
}

func TestApplyModule(t *testing.T) {
	fs := afero.NewMemMapFs()

	e := applyModule(fs, "mymodule", "../util/test-module", templates.Templates.ModuleInvocation)
	assert.Nil(t, e)

	s, e := fs.Stat("mymodule")
	assert.Nil(t, e)
	assert.True(t, s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	assert.Nil(t, e)
	r, e := afero.ReadFile(fs, "mymodule/main.tf")
	assert.Nil(t, e)
	expected := `module "test-module" {
  source = "../util/test-module"
  bar    = "${var.bar}"
  foo    = "${var.foo}"
}
`
	assert.Equal(t, expected, string(r))

	_, e = fs.Stat("mymodule/outputs.tf")
	assert.Nil(t, e)
	r, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	assert.Nil(t, e)
	expected = `output "bar" {
  value = "${module.test-module.bar}"
}

output "foo" {
  value = "${module.test-module.foo}"
}
`
	assert.Equal(t, expected, string(r))
}
func TestGetTargetPath(t *testing.T) {
	data := []struct {
		base    string
		source  string
		siccOff string
		siccOn  string
	}{
		{"", "foo.tmpl", "foo", "foo"},
		{"", "foo.tf.tmpl", "foo.tf", "foo.tf"},
		{"", "fogg.tf", "fogg.tf", "sicc.tf"},
		{"", "fogg.tf.tmpl", "fogg.tf", "sicc.tf"},
		{"foo", "foo.tmpl", "foo/foo", "foo/foo"},
		{"foo", "foo.tf.tmpl", "foo/foo.tf", "foo/foo.tf"},
		{"foo", "fogg.tf", "foo/fogg.tf", "foo/sicc.tf"},
		{"foo", "fogg.tf.tmpl", "foo/fogg.tf", "foo/sicc.tf"},
	}
	for _, test := range data {
		t.Run(test.source, func(t *testing.T) {
			off := getTargetPath(test.base, test.source, false)
			on := getTargetPath(test.base, test.source, true)
			assert.Equal(t, test.siccOff, off)
			assert.Equal(t, test.siccOn, on)

		})
	}
}

func TestFmtHcl(t *testing.T) {
	before := `foo { bar     = "bam"}`
	after := `foo {
  bar = "bam"
}
`
	fs := afero.NewMemMapFs()
	in := strings.NewReader(before)
	e := afero.WriteReader(fs, "foo.tf", in)
	assert.Nil(t, e)
	e = fmtHcl(fs, "foo.tf")
	assert.Nil(t, e)
	out, e := afero.ReadFile(fs, "foo.tf")
	assert.Nil(t, e)
	assert.NotNil(t, out)
	s := string(out)
	assert.Equal(t, after, s)
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
