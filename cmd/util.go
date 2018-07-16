package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	git "gopkg.in/src-d/go-git.v4"
)

func openGitOrExit(pwd string) *git.Repository {
	g, err := git.PlainOpen(pwd)
	if err != nil {
		// assuming this means no repository
		log.Fatal("fogg must be run from the root of a git repo")
		os.Exit(1)
	}
	return g
}

func readAndValidateConfig(fs afero.Fs, configFile string, verbose bool) (*config.Config, error) {
	config, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return nil, err
	}
	if verbose {
		log.Debug("CONFIG")
		log.Debug("%#v\n=====", config)
	}

	err = config.Validate()
	return config, err
}

func exitOnConfigErrors(err error) {
	if err != nil {
		log.Error("Config Error(s):")
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				log.Error("\t%s is a %s %s\n", err.Namespace(), err.Tag(), err.Kind())
			}
		} else {
			log.Panic(err)
		}
		os.Exit(1)
	}
}
