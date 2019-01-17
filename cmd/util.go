package cmd

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	validator "gopkg.in/go-playground/validator.v9"
)

func openGitOrExit(fs afero.Fs) {
	_, err := fs.Stat(".git")
	if err != nil {
		// assuming this means no repository
		log.Fatal("fogg must be run from the root of a git repo")
		os.Exit(1)
	}
}

func readAndValidateConfig(fs afero.Fs, configFile string, verbose bool) (*config.Config, error) {
	config, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return nil, errs.WrapUser(err, "unable to read config file")
	}
	if verbose {
		log.Debug("CONFIG")
		log.Debugf("%#v\n=====", config)
	}

	err = config.Validate()
	return config, err
}

func mergeConfigValidationErrors(err error) error {
	if err != nil {
		fmt.Println("fogg.json has error(s):")
		validatonErrors, ok := err.(validator.ValidationErrors)
		if ok {
			var sb strings.Builder
			for _, err := range validatonErrors {
				msg := fmt.Sprintf("\t%s is a %s %s\n", err.Namespace(), err.Tag(), err.Kind())
				sb.WriteString(strings.Replace(msg, "Config.", "", 1))
			}
			return errs.NewUser(sb.String())
		}
		return err
	}
	return nil
}

func setupDebug(debug bool) {
	logLevel := log.InfoLevel
	if debug { // debug overrides quiet
		logLevel = log.DebugLevel
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
			http.HandleFunc("/", pprof.Index)
		}()
	} else if quiet {
		logLevel = log.FatalLevel
	}
	log.SetLevel(logLevel)
}

func openFs() (afero.Fs, error) {
	pwd, e := os.Getwd()
	if e != nil {
		return nil, e
	}
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	return fs, nil
}

func bootstrapCmd(cmd *cobra.Command, debug bool) (afero.Fs, *config.Config, error) {
	setupDebug(debug)

	fs, err := openFs()
	if err != nil {
		return nil, nil, err
	}

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return nil, nil, err
	}

	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, nil, err
	}

	config, err := readAndValidateConfig(fs, configFile, verbose)

	err = mergeConfigValidationErrors(err)
	if err != nil {
		return nil, nil, err
	}

	return fs, config, nil
}
