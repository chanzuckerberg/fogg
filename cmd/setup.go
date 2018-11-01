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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	setupCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	setupCmd.Flags().BoolP("verbose", "v", false, "use this to turn on verbose output")
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

		// create the setup
		setup := setup{
			config: config,
			pwd:    pwd,
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

func (s *setup) getBinPath() string {
	return path.Join(s.pwd, plugins.BinDir, s.getOsArchPathComponent())
}

func (s *setup) getTfenvPath() string {
	return path.Join(s.getBinPath(), ".tfenv")
}

// tfEnv installs tfEnv
func (s *setup) tfEnv() error {
	// some paths
	tfenvPath := s.getTfenvPath()
	tfenvExecPathSrc := path.Join(tfenvPath, "bin", "tfenv")
	tfenvExecPathTarget := path.Join(s.getBinPath(), "tfenv")
	terraformExecPathSrc := path.Join(tfenvPath, "bin", "terraform")
	terraformExecPathTarget := path.Join(s.getBinPath(), "terraform")

	_, err := os.Stat(tfenvPath)
	if err == nil {
		// if something here, remove it to start clean
		// TODO: error handling?
		os.RemoveAll(tfenvPath)
		os.Remove(tfenvExecPathTarget)
		os.Remove(terraformExecPathTarget)
	}
	if err != nil && !os.IsNotExist(err) {
		return errs.WrapInternal(err, "Could not stat tfenv dir")
	}
	log.Debugf("Git clone to %s", tfenvPath)
	// not exist error
	cmd := exec.Command("git", "clone", "https://github.com/kamatama41/tfenv.git", tfenvPath)
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return errs.WrapInternal(err, "Could not clone tfenv")
	}

	// link tfenv and terraform
	err = os.Symlink(tfenvExecPathSrc, tfenvExecPathTarget)
	if err != nil {
		return errs.WrapInternal(err, "Could not link tfenv")
	}
	err = os.Symlink(terraformExecPathSrc, terraformExecPathTarget)
	if err != nil {
		return errs.WrapInternal(err, "Could not link terraform")
	}
	return nil
}

// terraform installs terraform
// assume tfenv already present
func (s *setup) terraform() error {
	pathEnv := fmt.Sprintf("PATH=%s:%s", os.Getenv("PATH"), s.getBinPath())
	log.Debugf("path: %s", pathEnv)
	cmd := exec.Command("tfenv", "install", s.config.Defaults.TerraformVersion)
	// cmd := exec.Command("printenv")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{pathEnv}

	err := cmd.Run()
	return errs.WrapInternal(err, "Could not install terraform")
}

// customProviders installs custom providers
func (s *setup) customProviders() error {
	return nil
}

// customPlugins installs custom plugins
func (s *setup) customPlugins() error {
	return nil
}
