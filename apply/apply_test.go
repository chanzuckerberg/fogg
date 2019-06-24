package apply

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	v1 "github.com/chanzuckerberg/fogg/config/v1"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	formatter := &logrus.TextFormatter{
		DisableTimestamp: true,
	}
	logrus.SetFormatter(formatter)
}

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func getNonExistentDirectoryName() string {
	nonexistentDir := "noexist-" + randomString(20)
	for {
		_, err := os.Stat(nonexistentDir)
		if os.IsNotExist(err) {
			return nonexistentDir
		}
		nonexistentDir = "noexist-" + randomString(20)
	}
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
	a := assert.New(t)
	sourceFile := strings.NewReader("foo")
	dest, d, err := util.TestFs()
	a.NoError(err)
	defer os.Remove(d)

	path := "bar"
	overrides := struct{ Foo string }{"foo"}

	e := applyTemplate(sourceFile, &templates.Templates.Common, dest, path, overrides)
	a.Nil(e)
	f, e := dest.Open("bar")
	a.Nil(e)
	r, e := ioutil.ReadAll(f)
	a.Nil(e)
	a.Equal("foo", string(r))
}

func TestApplyTemplateBasicNewDirectory(t *testing.T) {
	a := assert.New(t)
	sourceFile := strings.NewReader("foo")

	dest, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)

	nonexistentDir := getNonExistentDirectoryName()
	defer dest.RemoveAll(nonexistentDir) //nolint
	path := filepath.Join(nonexistentDir, "bar")
	overrides := struct{ Foo string }{"foo"}

	e := applyTemplate(sourceFile, &templates.Templates.Common, dest, path, overrides)
	a.Nil(e)
	f, e := dest.Open(path)
	a.Nil(e)
	r, e := ioutil.ReadAll(f)
	a.Nil(e)
	a.Equal("foo", string(r))
}

func TestApplyTemplate(t *testing.T) {
	a := assert.New(t)
	sourceFile := strings.NewReader("Hello {{.Name}}")
	dest, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)

	path := "hello"
	overrides := struct{ Name string }{"World"}

	e := applyTemplate(sourceFile, &templates.Templates.Common, dest, path, overrides)
	a.Nil(e)
	f, e := dest.Open("hello")
	a.Nil(e)
	r, e := ioutil.ReadAll(f)
	a.Nil(e)
	a.Equal("Hello World", string(r))
}

func TestTouchFile(t *testing.T) {
	a := assert.New(t)
	fs, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)

	e := touchFile(fs, "foo")
	a.Nil(e)
	r, e := readFile(fs, "foo")
	a.Nil(e)
	a.Equal("", r)

	fs, d2, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d2)

	err = writeFile(fs, "asdf", "jkl")
	a.NoError(err)

	r, e = readFile(fs, "asdf")
	a.Nil(e)
	a.Equal("jkl", r)

	e = touchFile(fs, "asdf")
	a.Nil(e)
	r, e = readFile(fs, "asdf")
	a.Nil(e)
	a.Equal("jkl", r)
}

func TestTouchFileNonExistentDirectory(t *testing.T) {
	a := assert.New(t)
	dest, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)

	nonexistentDir := getNonExistentDirectoryName()
	defer dest.RemoveAll(nonexistentDir) //nolint
	e := touchFile(dest, filepath.Join(nonexistentDir, "foo"))
	a.Nil(e)
	r, e := readFile(dest, filepath.Join(nonexistentDir, "foo"))
	a.Nil(e)
	a.Equal("", r)
	a.Nil(e)
}

func TestCreateFile(t *testing.T) {
	a := assert.New(t)
	fs, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)

	// create new file

	e := createFile(fs, "foo", strings.NewReader("bar"))
	a.Nil(e)

	r, e := readFile(fs, "foo")
	a.Nil(e)
	a.Equal("bar", r)

	// not overwrite existing file

	fs, d2, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d2)

	e = createFile(fs, "foo", strings.NewReader("bar"))
	a.Nil(e)

	r, e = readFile(fs, "foo")
	a.Nil(e)
	a.Equal("bar", r)

	e = createFile(fs, "foo", strings.NewReader("BAM"))
	a.Nil(e)

	r, e = readFile(fs, "foo")
	a.Nil(e)
	a.Equal("bar", r)
}

func TestCreateFileNonExistentDirectory(t *testing.T) {
	a := assert.New(t)
	// create new file in nonexistent directory
	dest, _, e := util.TestFs()
	a.NoError(e)

	e = createFile(dest, "newdir/foo", strings.NewReader("bar"))
	a.NoError(e)

	r, e := readFile(dest, "newdir/foo")
	a.NoError(e)
	a.Equal("bar", r)

}

func TestApplySmokeTest(t *testing.T) {
	r := require.New(t)
	fs, _, err := util.TestFs()
	r.NoError(err)

	yml := `
defaults:
  owner: foo
  project: bar
  terraform_version: 0
  backend:
    bucket: baz
    region: qux
    profile: quux
version: 2
`
	b, e := yaml.Marshal(yml)
	r.NoError(e)
	e = afero.WriteFile(fs, "fogg.yml", b, 0644)
	r.NoError(e)
	c, e := v2.ReadConfig([]byte(yml), fs, "fogg.yml")
	r.NoError(e)

	w, e := c.Validate()
	r.NoError(e)
	r.Len(w, 0)

	e = Apply(fs, c, templates.Templates, false)
	r.NoError(e)
}

func TestApplySmokeTestJSON(t *testing.T) {
	r := require.New(t)
	fs, _, err := util.TestFs()
	r.NoError(err)

	json := `
	{
	  "defaults": {
		"aws_region_provider": "reg",
		"aws_region_backend": "reg",
		"aws_profile_provider": "prof",
		"aws_profile_backend": "prof",
		"aws_provider_version": "0.12.0",
		"account_id": 789,
		"infra_s3_bucket": "buck",
		"project": "proj",
		"terraform_version": "0.100.0",
		"owner": "foo@example.com"
	  },
	  "travis_ci": {
		"enabled": true,
		"aws_iam_role_name": "travis",
			"id_account_name": "id",
			"test_buckets": 7
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
	c, e := v1.ReadConfig([]byte(json))
	r.NoError(e)
	c2, e := config.UpgradeConfigVersion(c)
	r.NoError(e)

	w, e := c2.Validate()
	r.NoError(e)
	r.Len(w, 1)

	e = Apply(fs, c2, templates.Templates, false)
	r.NoError(e)
}

func TestApplyModuleInvocation(t *testing.T) {
	a := assert.New(t)
	testFs, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)
	pwdFs, err := util.PwdFs()
	a.NoError(err)

	fs := afero.NewCopyOnWriteFs(pwdFs, testFs)

	e := applyModuleInvocation(fs, "mymodule", "test-module", templates.Templates.ModuleInvocation, &templates.Templates.Common)
	a.NoError(e)

	s, e := fs.Stat("mymodule")
	a.Nil(e)
	a.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	a.Nil(e)
	r, e := afero.ReadFile(fs, "mymodule/main.tf")
	a.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule \"test-module\" {\n  source = \"../test-module\"\n  bar    = \"${local.bar}\"\n  foo    = \"${local.foo}\"\n}\n"
	a.Equal(expected, string(r))

	_, e = fs.Stat("mymodule/outputs.tf")
	a.Nil(e)
	r, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	a.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\noutput \"bar\" {\n  value = \"${module.test-module.bar}\"\n}\n\noutput \"foo\" {\n  value = \"${module.test-module.foo}\"\n}\n"
	a.Equal(expected, string(r))
}
func TestGetTargetPath(t *testing.T) {
	data := []struct {
		base   string
		source string
		output string
	}{
		{"", "foo.tmpl", "foo"},
		{"", "foo.tf.tmpl", "foo.tf"},
		{"", "fogg.tf", "fogg.tf"},
		{"", "fogg.tf.tmpl", "fogg.tf"},
		{"foo", "foo.tmpl", "foo/foo"},
		{"foo", "foo.tf.tmpl", "foo/foo.tf"},
		{"foo", "fogg.tf", "foo/fogg.tf"},
		{"foo", "fogg.tf.tmpl", "foo/fogg.tf"},
	}
	for _, test := range data {
		t.Run(test.source, func(t *testing.T) {
			out := getTargetPath(test.base, test.source)
			assert.Equal(t, test.output, out)
		})
	}
}

func TestFmtHcl(t *testing.T) {
	a := assert.New(t)

	before := `foo { bar     = "bam"}`
	after := `foo {
  bar = "bam"
}
`
	fs, d, err := util.TestFs()
	a.NoError(err)
	defer os.RemoveAll(d)

	in := strings.NewReader(before)
	e := afero.WriteReader(fs, "foo.tf", in)
	a.Nil(e)
	e = fmtHcl(fs, "foo.tf")
	a.Nil(e)
	out, e := afero.ReadFile(fs, "foo.tf")
	a.Nil(e)
	a.NotNil(out)
	s := string(out)
	a.Equal(after, s)
}

func TestCalculateLocalPath(t *testing.T) {
	data := []struct {
		path          string
		moduleAddress string
		expected      string
	}{
		{"foo/bar", "bam/baz", "../../bam/baz"},
		{
			"foo/bar",
			"github.com/terraform-aws-modules/terraform-aws-vpc?ref=v1.30.0",
			"github.com/terraform-aws-modules/terraform-aws-vpc?ref=v1.30.0",
		},
		{"foo/bar", "github.com/asdf/jkl", "github.com/asdf/jkl"},
		// TODO modules from the registry don't work because it is
		// ambigious with file paths need to figure out how terraform
		// does this internally
		// {"foo/bar", "from/the/registry", "from/the/registry"},
	}

	for _, test := range data {
		t.Run("", func(t *testing.T) {
			p, e := calculateModuleAddressForSource(test.path, test.moduleAddress)
			assert.Nil(t, e)
			assert.Equal(t, test.expected, p)
		})
	}
}

var versionTests = []struct {
	current string
	tool    string
	result  bool
}{
	{"0.0.0", "0.1.0", true},
	{"0.1.0", "0.0.0", true},
	{"0.1.0", "0.1.0", false},
	{"0.1.0", "0.1.0-abcdef", true},
	{"0.1.0", "0.1.0-abcdef-dirty", true},
	{"0.1.0-abcdef-dirty", "0.1.0", true},
	{"0.1.0-abc", "0.1.0-def", true},
	{"0.1.0-abc", "0.1.0-abc", false},
}

func TestCheckToolVersions(t *testing.T) {
	for _, tc := range versionTests {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			fs, d, err := util.TestFs()
			a.NoError(err)
			defer os.RemoveAll(d)

			err = writeFile(fs, ".fogg-version", tc.current)
			a.NoError(err)

			v, _, e := checkToolVersions(fs, tc.tool)
			a.NoError(e)
			a.Equal(tc.result, v)
		})
	}
}

func TestVersionIsChanged(t *testing.T) {
	a := assert.New(t)

	for _, test := range versionTests {
		t.Run("", func(t *testing.T) {
			b, e := versionIsChanged(test.current, test.tool)
			a.NoError(e)
			a.Equal(test.result, b)
		})
	}
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

func Test_filepathRel(t *testing.T) {
	type args struct {
		name string
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"terraform.d", args{"terraform/accounts/idseq/terraform.d", "terraform.d"}, "../../../terraform.d", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filepathRel(tt.args.name, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("filepathRel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("filepathRel() = %v, want %v", got, tt.want)
			}
		})
	}
}
