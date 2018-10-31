package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	setupCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup dependencies for this project",
	Long: `This command will set up dependencies for this project.
				 These include things like tfenv, terraform, and custom plugins.`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		setupDebug(debug)
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapInternal(err, "Could not Getwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return errs.WrapInternal(err, "Could not parse verbose flag")
		}
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return errs.WrapInternal(err, "Could not parse config flag")
		}

		// parse config
		config, err := readAndValidateConfig(fs, configFile, verbose)
		err = mergeConfigValidationErrors(err)
		if err != nil {
			return err
		}

		// check that we are at root of initialized git repo
		openGitOrExit(pwd)
		setup := setup{
			config: config,
		}
		err = setup.tfEnv()
		if err != nil {
			return err
		}
		err = setup.terraform()
		if err != nil {
			return err
		}
		err = setup.customProviders()
		if err != nil {
			return err
		}
		return setup.customPlugins()
	},
}

// is a namespace
type setup struct {
	pwd    string
	config *config.Config
	fs     afero.Fs
}

func (s *setup) getOsArchPathComponent() string {
	return fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
}

func (s *setup) getTfenvPath() string {
	return path.Join(s.pwd, plugins.BinDir, s.getOsArchPathComponent(), ".tfenv")
}

// tfEnv installs tfEnv
func (s *setup) tfEnv() error {
	tfenvPath := s.getTfenvPath()
	_, err := os.Stat(tfenvPath)
	if err == nil {
		// if no error, then presumably we're done here
		return nil
	}
	if !os.IsNotExist(err) {
		return errs.WrapInternal(err, "Could not stat tfenv dir")
	}
	// not exist error
	cmd := exec.Command("git", "clone", "https://github.com/kamatama41/tfenv.git", tfenvPath)
	cmd.Env = os.Environ()

	err = cmd.Run()
	if err != nil {
		return errs.WrapInternal(err, "Could not clone tfenv")
	}

	return nil
}

// terraform installs terraform
func (s *setup) terraform() error {
	return nil
}

// customProviders installs custom providers
func (s *setup) customProviders() error {
	return nil
}

// customPlugins installs custom plugins
func (s *setup) customPlugins() error {
	return nil
}
