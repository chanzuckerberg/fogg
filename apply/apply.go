package apply

import (
	"log"
	"path/filepath"

	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/gobuffalo/packr"
	"github.com/spf13/afero"
)

func Apply(fs afero.Fs, configFile string, tmp *templates.T) error {
	p, err := plan.Eval(fs, configFile)
	// util.Dump(err)
	// util.Dump(p)

	applyRepo(fs, p, &tmp.Repo)
	return nil
}

func applyRepo(fs afero.Fs, p *plan.Plan, repoBox *packr.Box) error {
	applyTree(repoBox, fs, nil)
	return nil
}

func applyTree(source *packr.Box, dest afero.Fs, subst interface{}) error {
	source.Walk(func(path string, sourceFile packr.File) error {
		extension := filepath.Ext(path)
		// util.Dump(path)
		// util.Dump(sourceFile)
		// util.Dump(extension)
		if extension == ".tpl" {
			// if ext == '.tpl':
			//     dest, _ = os.path.splitext(dest)
			//     template(source, dest, substitutions)
			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])
			// elif ext == '.create':
			//     dest, _ = os.path.splitext(dest)
			//     if not os.path.exists(dest):
			//         shutil.copy(source, dest)
			//         if dest.endswith('.tf'):
			//             subprocess.call(['terraform', 'fmt', dest])
		} else if extension == ".touch" {
			log.Printf("touching %s", path)
			dest.Create(path)
			//     dest, _ = os.path.splitext(dest)
			//     print("touching %s" % dest)
			//     Path(dest).touch()
			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])

			// elif ext == '.rm':
			//     dest, _ = os.path.splitext(dest)
			//     print("removing %s" % dest)
			//     silentremove(dest)
		} else {
			log.Printf("copying %s", path)
			afero.WriteReader(dest, path, sourceFile)
		}
		return nil
	})
	return nil
}
