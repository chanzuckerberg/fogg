package apply

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/gobuffalo/packr"
	"github.com/spf13/afero"
)

func Apply(fs afero.Fs, configFile string, tmp *templates.T) error {
	p, err := plan.Eval(fs, configFile)
	util.Dump(err)
	util.Dump(p)

	applyRepo(fs, p, &tmp.Repo)
	return nil
}

func applyRepo(fs afero.Fs, p *plan.Plan, repoBox *packr.Box) error {
	applyTree(repoBox, fs, p)
	return nil
}

func applyTree(source *packr.Box, dest afero.Fs, subst interface{}) error {
	source.Walk(func(path string, sourceFile packr.File) error {
		extension := filepath.Ext(path)
		// util.Dump(path)
		// util.Dump(sourceFile)
		// util.Dump(extension)
		if extension == ".tmpl" {
			d := removeExtension(path)
			log.Printf("templating %s", d)
			writer, _ := dest.OpenFile(d, os.O_RDWR|os.O_CREATE, 0755)
			applyTemplate(sourceFile, writer, subst)

			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])
		} else if extension == ".touch" {
			d := removeExtension(path)
			_, err := dest.Stat(d)
			if err != nil { // TODO we might not want to do this for all errors
				log.Printf("touching %s", d)
				dest.Create(d)
			} else {
				log.Printf("skipping touch on existing file %s", d)
			}
			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])

		} else if extension == ".create" {
			d := removeExtension(path)
			_, err := dest.Stat(d)
			if err != nil { // TODO we might not want to do this for all errors
				log.Printf("creating %s", d)
				afero.WriteReader(dest, path, sourceFile)
			} else {
				log.Printf("skipping create on existing file %s", d)
			}
			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])

		} else {
			log.Printf("copying %s", path)
			afero.WriteReader(dest, path, sourceFile)
		}
		return nil
	})
	return nil
}

func removeExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))

}

func joinEnvs(m map[string]plan.Env) string {
	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, " ")
}

func applyTemplate(source packr.File, dest io.Writer, overrides interface{}) error {
	s, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}
	t := template.Must(template.New("tmpl").Funcs(sprig.FuncMap()).Parse(string(s)))
	return t.Execute(dest, overrides)

}
