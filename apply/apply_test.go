package apply

import (
	"encoding/json"
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
			planAccounts[local] = plan.Account{ComponentCommon: plan.ComponentCommon{
				Backend: plan.Backend{
					Kind: plan.BackendKindRemote,
				},
			}}
		}
		planEnvs := map[string]plan.Env{}
		for _, local := range test.planEnvs {
			splits := strings.Split(local, "/")
			if len(splits) != 2 {
				return nil, errors.New("env needs to be of the form env/component")
			}
			component := map[string]plan.Component{}
			component[splits[1]] = plan.Component{ComponentCommon: plan.ComponentCommon{
				Backend: plan.Backend{
					Kind: plan.BackendKindRemote,
				},
			}}
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
	downloader, err := util.MakeDownloader("test-module")
	r.NoError(err)
	comp := plan.Component{
		ModuleSource: util.StrPtr("test-module"),
	}
	e := applyModuleInvocation(fs, "mymodule", comp, templates.Templates.ModuleInvocation, templates.Templates.Common, downloader)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule \"test-module\" {\n  source = \"../test-module\"\n  bar    = local.bar\n  baz    = local.baz\n  foo    = local.foo\n  quux   = local.quux\n\n\n\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\noutput \"bar\" {\n  value     = module.test-module.bar\n  sensitive = false\n}\n\noutput \"foo\" {\n  value     = module.test-module.foo\n  sensitive = false\n}\n\n\n"
	r.Equal(expected, string(i))
}

func TestTFEConfigOmitEmpty(t *testing.T) {
	r := require.New(t)
	testFS, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)
	pwdFS, err := util.PwdFs()
	r.NoError(err)
	afs := afero.NewCopyOnWriteFs(pwdFS, testFS)

	plan := &plan.Plan{
		TFE: &plan.TFEConfig{},
		Envs: map[string]plan.Env{
			"prod": {
				Components: map[string]plan.Component{
					"test": {
						ComponentCommon: plan.ComponentCommon{
							Backend: plan.Backend{
								Kind: plan.BackendKindRemote,
							},
						}},
				},
				Env: "prod",
			},
		},
	}
	existingLocals := LocalsTFE{
		Locals: &Locals{
			Envs: map[string]map[string]*TFEWorkspace{
				"prod": {
					"test": &TFEWorkspace{},
				},
			},
		},
	}
	path := fmt.Sprintf("%s/tfe", rootPath)
	err = afs.MkdirAll(path, 0755)
	r.NoError(err)
	localsPath := filepath.Join("terraform", "tfe", "locals.tf.json")
	existingLocalsFile, err := afs.OpenFile(localsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	r.NoError(err)
	err = json.NewEncoder(existingLocalsFile).Encode(&existingLocals)
	r.NoError(err)
	existingLocalsFile.Close()

	err = applyTFE(afs, plan, templates.Templates)
	r.NoError(err)
	alteredLocalsFile, err := afs.Open(localsPath)
	r.NoError(err)
	b, err := io.ReadAll(alteredLocalsFile)
	r.NoError(err)
	r.Equal(
		strings.Join(strings.Fields(`{"locals":{"envs":{"prod":{"test":{}}}}}`), ""),
		strings.Join(strings.Fields(string(b)), ""))
}

func TestTFEConfig(t *testing.T) {
	r := require.New(t)
	testFS, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)
	pwdFS, err := util.PwdFs()
	r.NoError(err)
	afs := afero.NewCopyOnWriteFs(pwdFS, testFS)

	plan := &plan.Plan{
		TFE: &plan.TFEConfig{},
		Envs: map[string]plan.Env{
			"prod": {
				Components: map[string]plan.Component{
					"test": {
						ComponentCommon: plan.ComponentCommon{
							Backend: plan.Backend{
								Kind: plan.BackendKindRemote,
							},
						},
					},
				},
				Env: "prod",
			},
		},
	}
	emptyTrigger := []string{}
	tfVersion := "0.100.0"
	existingLocals := LocalsTFE{
		Locals: &Locals{
			Envs: map[string]map[string]*TFEWorkspace{
				"prod": {
					"test": &TFEWorkspace{
						TriggerPrefixes:  &emptyTrigger,
						TerraformVersion: &tfVersion,
					},
				},
			},
		},
	}
	path := fmt.Sprintf("%s/tfe", rootPath)
	err = afs.MkdirAll(path, 0755)
	r.NoError(err)
	localsPath := filepath.Join("terraform", "tfe", "locals.tf.json")
	existingLocalsFile, err := afs.OpenFile(localsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	r.NoError(err)
	err = json.NewEncoder(existingLocalsFile).Encode(&existingLocals)
	r.NoError(err)
	existingLocalsFile.Close()

	err = applyTFE(afs, plan, templates.Templates)
	r.NoError(err)
	locals := LocalsTFE{}
	alteredLocalsFile, err := afs.Open(localsPath)
	r.NoError(err)
	err = json.NewDecoder(alteredLocalsFile).Decode(&locals)
	r.NoError(err)
	_, ok := locals.Locals.Envs["prod"]
	r.True(ok)
	_, ok = locals.Locals.Envs["prod"]["test"]
	r.True(ok)
	r.Equal(*existingLocals.Locals.Envs["prod"]["test"].TerraformVersion, *locals.Locals.Envs["prod"]["test"].TerraformVersion)
	r.Equal(*existingLocals.Locals.Envs["prod"]["test"].TriggerPrefixes, *locals.Locals.Envs["prod"]["test"].TriggerPrefixes)
}

func TestApplyModuleInvocationWithEmptyVariables(t *testing.T) {
	r := require.New(t)
	testFs, d, err := util.TestFs()
	r.NoError(err)
	defer os.RemoveAll(d)
	pwdFs, err := util.PwdFs()
	r.NoError(err)

	fs := afero.NewCopyOnWriteFs(pwdFs, testFs)
	downloader, err := util.MakeDownloader("test-module")
	r.NoError(err)
	comp := plan.Component{ModuleSource: util.StrPtr("test-module"), Variables: []string{}}
	e := applyModuleInvocation(fs, "mymodule", comp, templates.Templates.ModuleInvocation, templates.Templates.Common, downloader)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule \"test-module\" {\n  source = \"../test-module\"\n  bar    = local.bar\n  foo    = local.foo\n\n\n\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\noutput \"bar\" {\n  value     = module.test-module.bar\n  sensitive = false\n}\n\noutput \"foo\" {\n  value     = module.test-module.foo\n  sensitive = false\n}\n\n\n"
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
	downloader, err := util.MakeDownloader("test-module")
	r.NoError(err)
	comp := plan.Component{
		ModuleSource: util.StrPtr("test-module"),
		Variables:    []string{"baz"},
	}
	e := applyModuleInvocation(fs, "mymodule", comp, templates.Templates.ModuleInvocation, templates.Templates.Common, downloader)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule \"test-module\" {\n  source = \"../test-module\"\n  bar    = local.bar\n  baz    = local.baz\n  foo    = local.foo\n\n\n\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\noutput \"bar\" {\n  value     = module.test-module.bar\n  sensitive = false\n}\n\noutput \"foo\" {\n  value     = module.test-module.foo\n  sensitive = false\n}\n\n\n"
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
	downloader, err := util.MakeDownloader("test-module")
	r.NoError(err)
	comp := plan.Component{
		ModuleSource: util.StrPtr("test-module"),
		ModuleName:   util.StrPtr(moduleName),
	}
	e := applyModuleInvocation(fs, "mymodule", comp, templates.Templates.ModuleInvocation, templates.Templates.Common, downloader)
	r.NoError(e)

	s, e := fs.Stat("mymodule")
	r.Nil(e)
	r.True(s.IsDir())

	_, e = fs.Stat("mymodule/main.tf")
	r.Nil(e)
	i, e := afero.ReadFile(fs, "mymodule/main.tf")
	r.Nil(e)
	expected := "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\nmodule \"module-name\" {\n  source = \"../test-module\"\n  bar    = local.bar\n  baz    = local.baz\n  foo    = local.foo\n  quux   = local.quux\n\n\n\n}\n"
	r.Equal(expected, string(i))

	_, e = fs.Stat("mymodule/outputs.tf")
	r.Nil(e)
	i, e = afero.ReadFile(fs, "mymodule/outputs.tf")
	r.Nil(e)
	expected = "# Auto-generated by fogg. Do not edit\n# Make improvements in fogg, so that everyone can benefit.\n\noutput \"bar\" {\n  value     = module.module-name.bar\n  sensitive = false\n}\n\noutput \"foo\" {\n  value     = module.module-name.foo\n  sensitive = false\n}\n\n\n"
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
