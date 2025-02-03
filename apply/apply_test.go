package apply

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/plan"
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

type UploadLocalsTestCase struct {
	locals                     LocalsTFE
	plan                       plan.Plan
	resultEnvs, resultAccounts []string
}

type UploadLocalsTestCaseSimple struct {
	localsAccounts, localEnvs  []string
	planAccounts, planEnvs     []string
	resultEnvs, resultAccounts []string
}

func makeTestCases(tests []UploadLocalsTestCaseSimple) ([]UploadLocalsTestCase, error) {
	testcases := []UploadLocalsTestCase{}
	for _, test := range tests {
		testcase := UploadLocalsTestCase{}

		localAccounts := make(map[string]*TFEWorkspace, 0)
		for _, local := range test.localsAccounts {
			localAccounts[local] = MakeTFEWorkspace("1.2.6")
		}
		localEnvs := make(map[string]map[string]*TFEWorkspace, 0)
		for _, local := range test.localEnvs {
			splits := strings.Split(local, "/")
			if len(splits) != 2 {
				return nil, errors.New("env needs to be of the form env/component")
			}
			component := map[string]*TFEWorkspace{}
			component[splits[1]] = MakeTFEWorkspace("1.2.6")
			localEnvs[splits[0]] = component
		}
		locals := LocalsTFE{
			Locals: &Locals{
				Accounts: localAccounts,
				Envs:     localEnvs,
			},
		}

		planAccounts := map[string]plan.Account{}
		for _, local := range test.planAccounts {
			planAccounts[local] = plan.Account{}
		}
		planEnvs := map[string]plan.Env{}
		for _, local := range test.planEnvs {
			splits := strings.Split(local, "/")
			if len(splits) != 2 {
				return nil, errors.New("env needs to be of the form env/component")
			}
			component := map[string]plan.Component{}
			component[splits[1]] = plan.Component{}
			planEnvs[splits[0]] = plan.Env{
				Components: component,
			}
		}
		plan := plan.Plan{
			Accounts: planAccounts,
			Envs:     planEnvs,
		}
		testcase.locals = locals
		testcase.plan = plan
		testcase.resultAccounts = test.resultAccounts
		testcase.resultEnvs = test.resultEnvs
		testcases = append(testcases, testcase)
	}
	return testcases, nil
}

func TestUpdateLocalsFromPlan(t *testing.T) {
	r := require.New(t)
	testcases, err := makeTestCases([]UploadLocalsTestCaseSimple{
		{
			localsAccounts: []string{"a", "b"},
			localEnvs:      []string{"playground/a", "dev/b"},
			planAccounts:   []string{"a", "b"},
			planEnvs:       []string{"playground/a", "dev/b"},
			resultAccounts: []string{"a", "b"},
			resultEnvs:     []string{"playground/a", "dev/b"},
		},
		{
			localsAccounts: []string{"b"},
			localEnvs:      []string{"playground/a", "dev/b"},
			planAccounts:   []string{"a", "b"},
			planEnvs:       []string{"playground/a", "dev/b"},
			resultAccounts: []string{"a", "b"},
			resultEnvs:     []string{"playground/a", "dev/b"},
		},
		{
			localsAccounts: []string{},
			localEnvs:      []string{"playground/a", "dev/b"},
			planAccounts:   []string{"a", "b"},
			planEnvs:       []string{"playground/a", "dev/b"},
			resultAccounts: []string{"a", "b"},
			resultEnvs:     []string{"playground/a", "dev/b"},
		},
		{
			localsAccounts: []string{},
			localEnvs:      []string{},
			planAccounts:   []string{"a", "b"},
			planEnvs:       []string{"playground/a", "dev/b"},
			resultAccounts: []string{"a", "b"},
			resultEnvs:     []string{"playground/a", "dev/b"},
		},
		{
			localsAccounts: []string{"a", "b"},
			localEnvs:      []string{"playground/a", "dev/b"},
			planAccounts:   []string{"b"},
			planEnvs:       []string{"playground/a", "dev/b"},
			resultAccounts: []string{"b"},
			resultEnvs:     []string{"playground/a", "dev/b"},
		},
		{
			localsAccounts: []string{"a", "b"},
			localEnvs:      []string{"playground/a", "dev/b"},
			planAccounts:   []string{"b"},
			planEnvs:       []string{"dev/b"},
			resultAccounts: []string{"b"},
			resultEnvs:     []string{"dev/b"},
		},
		{
			localsAccounts: []string{"a", "b"},
			localEnvs:      []string{"playground/a", "dev/b"},
			planAccounts:   []string{"b"},
			planEnvs:       []string{"dev2/b"},
			resultAccounts: []string{"b"},
			resultEnvs:     []string{"dev2/b"},
		},
	})
	r.NoError(err)
	for _, testcase := range testcases {
		updateLocalsFromPlan(&testcase.locals, &testcase.plan)
		resultsAccounts := []string{}
		for key := range testcase.locals.Locals.Accounts {
			resultsAccounts = append(resultsAccounts, key)
		}
		resultEnvs := []string{}
		for env, envLocal := range testcase.locals.Locals.Envs {
			for component := range envLocal {
				resultEnvs = append(resultEnvs, fmt.Sprintf("%s/%s", env, component))
			}
		}
		r.ElementsMatch(testcase.resultAccounts, resultsAccounts)
		r.ElementsMatch(testcase.resultEnvs, resultEnvs)
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
	i, e := io.ReadAll(f)
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
	i, e := io.ReadAll(f)
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
	i, e := io.ReadAll(f)
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

	w, e := c.Validate(fs)
	r.NoError(e)
	r.Len(w, 0)

	e = Apply(fs, c, templates.Templates, false, nil, nil)
	r.NoError(e)
}

func TestApplyModuleInvocation(t *testing.T) {
	r := require.New(t)
	testFs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)
	pwdFs, err := util.PwdFs()
	r.NoError(err)

	moduleSource := "test-module"
	fs := afero.NewCopyOnWriteFs(pwdFs, testFs)
	downloader, err := util.MakeDownloader("test-module", "", nil)
	r.NoError(err)
	mi := []moduleInvocation{
		{
			module: v2.ComponentModule{
				Source:    &moduleSource,
				Variables: nil,
			},
			downloadFunc: downloader,
		},
	}
	_, e := applyModuleInvocation(fs, "mymodule", templates.Templates.ModuleInvocation, templates.Templates.Common, mi, nil)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule \"test-module\" {\n  source = \"../test-module\"\n  bar    = local.bar\n  baz    = local.baz\n  foo    = local.foo\n  quux   = local.quux\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\n# module \"test-module\" outputs\noutput \"bar\" {\n  value     = module.test-module.bar\n  sensitive = false\n}\noutput \"foo\" {\n  value     = module.test-module.foo\n  sensitive = false\n}\n"
	r.Equal(expected, string(i))
}

func TestApplyModuleInvocationWithEmptyVariables(t *testing.T) {
	r := require.New(t)
	testFs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)
	pwdFs, err := util.PwdFs()
	r.NoError(err)

	moduleSource := "test-module"
	fs := afero.NewCopyOnWriteFs(pwdFs, testFs)
	downloader, err := util.MakeDownloader(moduleSource, "", nil)
	r.NoError(err)
	mi := []moduleInvocation{
		{
			module: v2.ComponentModule{
				Source:    &moduleSource,
				Variables: []string{},
			},
			downloadFunc: downloader,
		},
	}
	_, e := applyModuleInvocation(fs, "mymodule", templates.Templates.ModuleInvocation, templates.Templates.Common, mi, nil)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule \"test-module\" {\n  source = \"../test-module\"\n  bar    = local.bar\n  foo    = local.foo\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\n# module \"test-module\" outputs\noutput \"bar\" {\n  value     = module.test-module.bar\n  sensitive = false\n}\noutput \"foo\" {\n  value     = module.test-module.foo\n  sensitive = false\n}\n"
	r.Equal(expected, string(i))
}

func TestApplyModuleInvocationWithOneDefaultVariable(t *testing.T) {
	r := require.New(t)
	testFs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)
	pwdFs, err := util.PwdFs()
	r.NoError(err)

	fs := afero.NewCopyOnWriteFs(pwdFs, testFs)
	downloader, err := util.MakeDownloader("test-module", "", nil)
	r.NoError(err)
	moduleName := "test-module"
	mi := []moduleInvocation{
		{
			module: v2.ComponentModule{
				Source:    &moduleName,
				Variables: []string{"baz"},
			},
			downloadFunc: downloader,
		},
	}
	_, e := applyModuleInvocation(fs, "mymodule", templates.Templates.ModuleInvocation, templates.Templates.Common, mi, nil)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule \"test-module\" {\n  source = \"../test-module\"\n  bar    = local.bar\n  baz    = local.baz\n  foo    = local.foo\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\n# module \"test-module\" outputs\noutput \"bar\" {\n  value     = module.test-module.bar\n  sensitive = false\n}\noutput \"foo\" {\n  value     = module.test-module.foo\n  sensitive = false\n}\n"
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

	moduleSource := "test-module"
	downloader, err := util.MakeDownloader(moduleSource, "", nil)
	r.NoError(err)
	moduleName := "module_name"
	mi := []moduleInvocation{
		{
			module: v2.ComponentModule{
				Name:      &moduleName,
				Source:    &moduleSource,
				Variables: nil,
			},
			downloadFunc: downloader,
		},
	}
	_, e := applyModuleInvocation(fs, "mymodule", templates.Templates.ModuleInvocation, templates.Templates.Common, mi, nil)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule \"module_name\" {\n  source = \"../test-module\"\n  bar    = local.bar\n  baz    = local.baz\n  foo    = local.foo\n  quux   = local.quux\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\n# module \"module_name\" outputs\noutput \"bar\" {\n  value     = module.module_name.bar\n  sensitive = false\n}\noutput \"foo\" {\n  value     = module.module_name.foo\n  sensitive = false\n}\n"
	r.Equal(expected, string(i))
}

func TestApplyModuleInvocationWithModulePrefix(t *testing.T) {
	r := require.New(t)
	testFs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)
	pwdFs, err := util.PwdFs()
	r.NoError(err)

	fs := afero.NewCopyOnWriteFs(pwdFs, testFs)

	downloader, err := util.MakeDownloader("test-module", "", nil)
	r.NoError(err)
	moduleName := "module_name"
	modulePrefix := "prefix"
	moduleSource := "test-module"
	mi := []moduleInvocation{
		{
			module: v2.ComponentModule{
				Name:      &moduleName,
				Prefix:    &modulePrefix,
				Source:    &moduleSource,
				Variables: []string{},
			},
			downloadFunc: downloader,
		},
	}
	_, e := applyModuleInvocation(fs, "mymodule", templates.Templates.ModuleInvocation, templates.Templates.Common, mi, nil)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule \"module_name\" {\n  source = \"../test-module\"\n  bar    = local.prefix_bar\n  foo    = local.prefix_foo\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\n# module \"module_name\" outputs\noutput \"prefix_bar\" {\n  value     = module.module_name.bar\n  sensitive = false\n}\noutput \"prefix_foo\" {\n  value     = module.module_name.foo\n  sensitive = false\n}\n"
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
			p, e := calculateModuleAddressForSource(tt.path, tt.moduleAddress, "")
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
	r, e := io.ReadAll(f)
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
func Test_dropDirectorySuffix(t *testing.T) {
	tests := []struct {
		path      string
		separator string
		expected  string
	}{
		{
			path:      "terraform/foo/bar",
			separator: "/",
			expected:  "terraform/foo/bar",
		},
		{
			path:      "terraform/foo.sample/bar",
			separator: "/",
			expected:  "terraform/foo/bar",
		},
		{
			path:      "terraform/foo.sample/bar.sample/baz.sample/qux",
			separator: "/",
			expected:  "terraform/foo/bar/baz/qux",
		},
		{
			path:      "terraform/foo.sample/bar.sample/baz.sample/qux.sample",
			separator: "/",
			expected:  "terraform/foo/bar/baz/qux",
		},
		{
			path:      `terraform\foo.sample\bar`,
			separator: `\`,
			expected:  `terraform\foo\bar`,
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			r := require.New(t)
			result := dropDirectorySuffix(test.path, test.separator)
			r.Equal(test.expected, result)
		})
	}
}

func TestApplyEnvsWithFilter(t *testing.T) {
	r := require.New(t)
	tmpl := templates.Templates
	dest, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)

	conf := &v2.Config{
		Defaults: v2.Defaults{
			Common: v2.Common{
				TerraformVersion: util.Ptr("0.12.0"),
				Owner:            util.Ptr("owner"),
				Project:          util.Ptr("project"),
				Backend: &v2.Backend{
					Kind:    util.Ptr(string(plan.BackendKindS3)),
					Bucket:  util.Ptr("bucket"),
					Region:  util.Ptr("region"),
					Profile: util.Ptr("profile"),
				},
			},
		},
		Envs: map[string]v2.Env{
			"env1": {
				Components: map[string]v2.Component{
					"component1": {},
				},
			},
			"env2": {
				Components: map[string]v2.Component{
					"component2": {},
				},
			},
		},
	}

	plan, err := plan.Eval(conf)
	r.NoError(err)
	envFilter := "env1"

	pathModuleConfigs, err := applyEnvs(dest, plan, &envFilter, nil, tmpl.Env, tmpl.Components, tmpl.Common, tmpl.TurboComponent)
	r.NoError(err)
	r.NotNil(pathModuleConfigs)
	_, err = dest.Stat("terraform/envs/env1")
	r.NoError(err)
	_, err = dest.Stat("terraform/envs/env2")
	r.Error(err)
}

func TestApplyEnvsWithCompFilter(t *testing.T) {
	r := require.New(t)
	tmpl := templates.Templates

	conf := &v2.Config{
		Defaults: v2.Defaults{
			Common: v2.Common{
				TerraformVersion: util.Ptr("0.12.0"),
				Owner:            util.Ptr("owner"),
				Project:          util.Ptr("project"),
				Backend: &v2.Backend{
					Kind:    util.Ptr(string(plan.BackendKindS3)),
					Bucket:  util.Ptr("bucket"),
					Region:  util.Ptr("region"),
					Profile: util.Ptr("profile"),
				},
			},
		},
		Envs: map[string]v2.Env{
			"env1": {
				Components: map[string]v2.Component{
					"component1": {},
					"component2": {},
				},
			},
			"env2": {
				Components: map[string]v2.Component{
					"component2": {},
				},
			},
		},
	}

	plan, err := plan.Eval(conf)
	r.NoError(err)

	tests := []struct {
		envFilter      *string
		compFilter     *string
		expectedDirs   []string
		unexpectedDirs []string
	}{
		{util.Ptr("env1"), util.Ptr("component1"), []string{"terraform/envs/env1/component1"}, []string{"terraform/envs/env1/component2", "terraform/envs/env2"}},
		{util.Ptr("env1"), util.Ptr("component2"), []string{"terraform/envs/env1/component2"}, []string{"terraform/envs/env1/component1", "terraform/envs/env2"}},
		{util.Ptr("env2"), util.Ptr("component2"), []string{"terraform/envs/env2/component2"}, []string{"terraform/envs/env1"}},
		{util.Ptr("env1"), nil, []string{"terraform/envs/env1/component1", "terraform/envs/env1/component2"}, []string{"terraform/envs/env2"}},
		{nil, nil, []string{"terraform/envs/env1/component1", "terraform/envs/env1/component2", "terraform/envs/env2/component2"}, []string{}},
	}

	for _, tt := range tests {
		var testName string
		if tt.envFilter == nil {
			testName = "all"
		} else {
			testName = *tt.envFilter
			if tt.compFilter != nil {
				testName += "_" + *tt.compFilter
			} else {
				testName += "_all"
			}
		}

		t.Run(testName, func(t *testing.T) {
			dest, d, err := util.TestFs()
			r.NoError(err)
			defer os.RemoveAll(d)

			pathModuleConfigs, err := applyEnvs(dest, plan, tt.envFilter, tt.compFilter, tmpl.Env, tmpl.Components, tmpl.Common, tmpl.TurboComponent)
			r.NoError(err)
			r.NotNil(pathModuleConfigs)

			for _, dir := range tt.expectedDirs {
				_, err = dest.Stat(dir)
				r.NoError(err)
			}

			for _, dir := range tt.unexpectedDirs {
				_, err = dest.Stat(dir)
				r.Error(err)
			}
		})
	}
}

func TestPreservePackageVersion(t *testing.T) {
	r := require.New(t)

	moduleName := "cdktf-test"
	defaultVersion := "0.0.0"
	modulePlan := plan.Module{
		Name:                  moduleName,
		Kind:                  util.Ptr(v2.ModuleKindCDKTF),
		Version:               defaultVersion,
		Publish:               false,
		Author:                "author",
		CdktfDependencies:     map[string]string{},
		CdktfDevDependencies:  map[string]string{},
		CdktfPeerDependencies: map[string]string{},
	}

	tests := []struct {
		version *string
	}{
		{nil},
		{util.Ptr("1.2.3")},
	}

	for _, tt := range tests {
		var expectedVersion string
		if tt.version == nil {
			expectedVersion = defaultVersion
		} else {
			expectedVersion = *tt.version
		}

		t.Run(expectedVersion, func(t *testing.T) {
			fs, d, err := util.TestFs()
			r.NoError(err)
			defer os.RemoveAll(d)
			path := fmt.Sprintf("%s/modules/%s", util.RootPath, moduleName)
			if tt.version != nil {
				writeMockPackageJson(fs, path, expectedVersion)
			}
			preservePackageJsonVersion(fs, &modulePlan, path+"/package.json")
			r.Equal(expectedVersion, modulePlan.Version)
		})
	}
}

func writeMockPackageJson(fs afero.Fs, path string, version string) error {
	packageJsonPath := path + "/package.json"
	e := fs.MkdirAll(path, 0755)
	if e != nil {
		return e
	}
	contents := fmt.Sprintf(`{ "version": "%s" }`, version)
	e = writeFile(fs, packageJsonPath, contents)
	return e
}
