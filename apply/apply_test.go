package apply

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"
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
	r := require.New(t)
	x := removeExtension("foo")
	r.Equal("foo", x)

	x = removeExtension("foo.txt")
	r.Equal("foo", x)

	x = removeExtension("foo.txt.asdf")
	r.Equal("foo.txt", x)

	x = removeExtension("foo/bar.txt")
	r.Equal("foo/bar", x)
}

func TestApplyTemplateBasic(t *testing.T) {
	r := require.New(t)
	sourceFile := strings.NewReader("foo")
	dest, d, err := util.TestFs()
	r.NoError(err)
	defer os.Remove(d)

	path := "bar"
	overrides := struct{ Foo string }{"foo"}

	e := applyTemplate(sourceFile, templates.Templates.Common, dest, path, overrides)
	r.Nil(e)
	f, e := dest.Open("bar")
	r.Nil(e)
	i, e := ioutil.ReadAll(f)
	r.Nil(e)
	r.Equal("foo", string(i))
}

func TestApplyTemplateBasicNewDirectory(t *testing.T) {
	r := require.New(t)
	sourceFile := strings.NewReader("foo")

	dest, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	nonexistentDir := getNonExistentDirectoryName()
	defer dest.RemoveAll(nonexistentDir) //nolint
	path := filepath.Join(nonexistentDir, "bar")
	overrides := struct{ Foo string }{"foo"}

	e := applyTemplate(sourceFile, templates.Templates.Common, dest, path, overrides)
	r.Nil(e)
	f, e := dest.Open(path)
	r.Nil(e)
	i, e := ioutil.ReadAll(f)
	r.Nil(e)
	r.Equal("foo", string(i))
}

func TestApplyTemplate(t *testing.T) {
	r := require.New(t)
	sourceFile := strings.NewReader("Hello {{.Name}}")
	dest, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	path := "hello"
	overrides := struct{ Name string }{"World"}

	e := applyTemplate(sourceFile, templates.Templates.Common, dest, path, overrides)
	r.Nil(e)
	f, e := dest.Open("hello")
	r.Nil(e)
	i, e := ioutil.ReadAll(f)
	r.Nil(e)
	r.Equal("Hello World", string(i))
}

func TestTouchFile(t *testing.T) {
	r := require.New(t)
	fs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	e := touchFile(fs, "foo")
	r.Nil(e)
	i, e := readFile(fs, "foo")
	r.Nil(e)
	r.Equal("", i)

	fs, d2, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d2)

	err = writeFile(fs, "asdf", "jkl")
	r.NoError(err)

	i, e = readFile(fs, "asdf")
	r.Nil(e)
	r.Equal("jkl", i)

	e = touchFile(fs, "asdf")
	r.Nil(e)
	i, e = readFile(fs, "asdf")
	r.Nil(e)
	r.Equal("jkl", i)
}

func TestTouchFileNonExistentDirectory(t *testing.T) {
	r := require.New(t)
	dest, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	nonexistentDir := getNonExistentDirectoryName()
	defer dest.RemoveAll(nonexistentDir) //nolint
	e := touchFile(dest, filepath.Join(nonexistentDir, "foo"))
	r.Nil(e)
	i, e := readFile(dest, filepath.Join(nonexistentDir, "foo"))
	r.Nil(e)
	r.Equal("", i)
	r.Nil(e)
}

func TestCreateFile(t *testing.T) {
	r := require.New(t)
	fs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	// create new file

	e := createFile(fs, "foo", strings.NewReader("bar"))
	r.Nil(e)

	i, e := readFile(fs, "foo")
	r.Nil(e)
	r.Equal("bar", i)

	// not overwrite existing file

	fs, d2, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d2)

	e = createFile(fs, "foo", strings.NewReader("bar"))
	r.Nil(e)

	i, e = readFile(fs, "foo")
	r.Nil(e)
	r.Equal("bar", i)

	e = createFile(fs, "foo", strings.NewReader("BAM"))
	r.Nil(e)

	i, e = readFile(fs, "foo")
	r.Nil(e)
	r.Equal("bar", i)
}

func TestCreateFileNonExistentDirectory(t *testing.T) {
	r := require.New(t)
	// create new file in nonexistent directory
	dest, _, e := util.TestFs()
	r.NoError(e)

	e = createFile(dest, "newdir/foo", strings.NewReader("bar"))
	r.NoError(e)

	i, e := readFile(dest, "newdir/foo")
	r.NoError(e)
	r.Equal("bar", i)
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
	c, e := v2.ReadConfig(fs, []byte(yml), "fogg.yml")
	r.NoError(e)

	w, e := c.Validate()
	r.NoError(e)
	r.Len(w, 0)

	e = Apply(fs, c, templates.Templates, false)
	r.NoError(e)
}

func TestApplyModuleInvocation(t *testing.T) {
	r := require.New(t)
	testFs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)
	pwdFs, err := util.PwdFs()
	r.NoError(err)

	fs := afero.NewCopyOnWriteFs(pwdFs, testFs)

	e := applyModuleInvocation(fs, "mymodule", "test-module", nil, templates.Templates.ModuleInvocation, templates.Templates.Common)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule test-module {\n  source = \"../test-module\"\n  bar    = local.bar\n  foo    = local.foo\n\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\noutput bar {\n  value = module.test-module.bar\n}\n\noutput foo {\n  value = module.test-module.foo\n}\n\n\n"
	r.Equal(expected, string(i))
}

func TestApplyModuleInvocationWithModuleName(t *testing.T) {
	r := require.New(t)
	testFs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)
	pwdFs, err := util.PwdFs()
	r.NoError(err)

	fs := afero.NewCopyOnWriteFs(pwdFs, testFs)

	moduleName := "module-name"
	e := applyModuleInvocation(fs, "mymodule", "test-module", &moduleName, templates.Templates.ModuleInvocation, templates.Templates.Common)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule module-name {\n  source = \"../test-module\"\n  bar    = local.bar\n  foo    = local.foo\n\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\noutput bar {\n  value = module.module-name.bar\n}\n\noutput foo {\n  value = module.module-name.foo\n}\n\n\n"
	r.Equal(expected, string(i))
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
		tt := test

		t.Run(test.source, func(t *testing.T) {
			r := require.New(t)
			out := getTargetPath(tt.base, tt.source)
			r.Equal(tt.output, out)
		})
	}
}

func TestFmtHcl(t *testing.T) {
	r := require.New(t)

	before := `foo { bar     = "bam"}`
	after := `foo { bar = "bam" }`
	fs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	in := strings.NewReader(before)
	e := afero.WriteReader(fs, "foo.tf", in)
	r.Nil(e)
	e = fmtHcl(fs, "foo.tf", false)
	r.Nil(e)
	out, e := afero.ReadFile(fs, "foo.tf")
	r.Nil(e)
	r.NotNil(out)
	s := string(out)
	r.Equal(after, s)
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
		tt := test

		t.Run("", func(t *testing.T) {
			r := require.New(t)
			p, e := calculateModuleAddressForSource(tt.path, tt.moduleAddress)
			r.Nil(e)
			r.Equal(tt.expected, p)
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
	{"\t\n 0.2.0 \n\t", "0.2.0", false},
}

func TestCheckToolVersions(t *testing.T) {
	for _, test := range versionTests {
		tt := test

		t.Run("", func(t *testing.T) {
			r := require.New(t)
			fs, d, err := util.TestFs()
			r.NoError(err)
			defer os.RemoveAll(d)

			err = writeFile(fs, ".fogg-version", tt.current)
			r.NoError(err)

			v, _, e := checkToolVersions(fs, tt.tool)
			r.NoError(e)
			r.Equal(tt.result, v)
		})
	}
}

func TestVersionIsChanged(t *testing.T) {
	r := require.New(t)

	for _, test := range versionTests {
		tt := test

		t.Run("", func(t *testing.T) {
			b := versionIsChanged(tt.current, tt.tool)
			r.Equal(tt.result, b)
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
		name string
		args args
		want string
	}{
		{"terraform.d", args{"terraform/accounts/idseq/terraform.d", "terraform.d"}, "../../../terraform.d"},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			got := filepathRel(tt.args.name, tt.args.path)
			if got != tt.want {
				t.Errorf("filepathRel() = %v, want %v", got, tt.want)
			}
		})
	}
}
