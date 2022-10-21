package cmd

import (
	"io"
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(fmtCmd)
}

var fmtCmd = &cobra.Command{
	Use:           "fmt",
	Short:         "format",
	Long:          "",
	SilenceErrors: true, // If we don't silence here, cobra will print them. But we want to do that in cmd/root.go
	RunE: func(cmd *cobra.Command, args []string) error {
		fs, err := pwdFs()
		if err != nil {
			return err
		}

		// if fogg.yml exists, read and format it
		if fileExists("fogg.yml") {
			n := &yaml.Node{}

			f, err := fs.Open("fogg.yml")
			if err != nil {
				return errs.NewUser("could not open fogg.yml")
			}
			defer f.Close()

			b, err := io.ReadAll(f)
			if err != nil {
				return errs.WrapUser(err, "could not read fogg.yml")
			}

			err = yaml.Unmarshal(b, n)
			if err != nil {
				return errs.WrapUser(err, "unable to parse yaml")
			}

			f, e := fs.OpenFile("fogg.yml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if e != nil {
				return errors.Wrap(e, "Unable to open fogg.yml for writing")
			}
			enc := yaml.NewEncoder(f)
			enc.SetIndent(2)
			return enc.Encode(n)
		}

		return nil
	},
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func pwdFs() (afero.Fs, error) {
	// Set up fs
	pwd, e := os.Getwd()
	if e != nil {
		return nil, errs.WrapUser(e, "can't get pwd")
	}
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	return fs, nil
}
