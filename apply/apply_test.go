package apply

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/templates"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
	formatter := &log.TextFormatter{
		DisableTimestamp: true,
	}
	log.SetFormatter(formatter)
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

func TestApplyTemplateBasicNewDirectory(t *testing.T) {
	sourceFile := strings.NewReader("foo")
	// Potential errors do not show up if using NewMemMapFs; needs real OS fs.
	dest := afero.NewOsFs()
	nonexistentDir := getNonExistentDirectoryName()
	defer dest.RemoveAll(nonexistentDir)
	path := filepath.Join(nonexistentDir, "bar")
	overrides := struct{ Foo string }{"foo"}

	e := applyTemplate(sourceFile, dest, path, overrides)
	assert.Nil(t, e)
	f, e := dest.Open(path)
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

func TestTouchFileNonExistentDirectory(t *testing.T) {
	// Potential errors do not show up if using NewMemMapFs; needs real OS fs.
	dest := afero.NewOsFs()
	nonexistentDir := getNonExistentDirectoryName()
	defer dest.RemoveAll(nonexistentDir)
	e := touchFile(dest, filepath.Join(nonexistentDir, "foo"))
	assert.Nil(t, e)
	r, e := readFile(dest, filepath.Join(nonexistentDir, "foo"))
	assert.Nil(t, e)
	assert.Equal(t, "", r)
	assert.Nil(t, e)
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

func TestCreateFileNonExistentDirectory(t *testing.T) {

	// create new file in nonexistent directory
	dest := afero.NewOsFs()

	e := createFile(dest, "newdir/foo", strings.NewReader("bar"))
	assert.Nil(t, e)

	r, e := readFile(dest, "newdir/foo")
	assert.Nil(t, e)
	assert.Equal(t, "bar", r)

}

func TestApplySmokeTest(t *testing.T) {
	t.Skip("doesn't currently work because afero doesn't support symlinks")
	// We have to use a BasePathFs so that we can calculate `RealPath` for symlinking. Afero doesn't support symlinks
	fs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")
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
	c, e := v1.ReadConfig(ioutil.NopCloser(strings.NewReader(json)))
	assert.NoError(t, e)
	c2, e := config.UpgradeConfigVersion(c)
	assert.NoError(t, e)

	e = Apply(fs, c2, templates.Templates, false)
	assert.NoError(t, e)
}

func TestApplyModuleInvocation(t *testing.T) {
	fs := afero.NewMemMapFs()

	e := applyModuleInvocation(fs, "mymodule", "../util/test-module", templates.Templates.ModuleInvocation)
	assert.Nil(t, e)

	s, e := fs.Stat("mymodule")
	assert.Nil(t, e)
	assert.True(t, s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	assert.Nil(t, e)
	r, e := afero.ReadFile(fs, "mymodule/main.tf")
	assert.Nil(t, e)
	expected := `# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

module "test-module" {
  source = "../../util/test-module"
  bar    = "${local.bar}"
  foo    = "${local.foo}"
}
`
	assert.Equal(t, expected, string(r))

	_, e = fs.Stat("mymodule/outputs.tf")
	assert.Nil(t, e)
	r, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	assert.Nil(t, e)
	expected = `# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

output "bar" {
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
	a := assert.New(t)

	for _, tc := range versionTests {
		t.Run("", func(t *testing.T) {
			fs := afero.NewMemMapFs()
			writeFile(fs, ".fogg-version", tc.current)

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
