package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	validator "gopkg.in/go-playground/validator.v9"
	git "gopkg.in/src-d/go-git.v4"
)

func openGitOrExit(pwd string) *git.Repository {
	log.Debugf("opening git at %s", pwd)
	g, err := git.PlainOpen(pwd)
	if err != nil {
		// assuming this means no repository
		log.Debug(errors.Wrap(err, "unable to open git index"))
		log.Fatal("fogg must be run from the root of a git repo")
		os.Exit(1)
	}
	return g
}

func readAndValidateConfig(fs afero.Fs, configFile string, verbose bool) (*config.Config, error) {
	config, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read config file")
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
		fmt.Println("fogg.json has error(s):")
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				msg := fmt.Sprintf("\t%s is a %s %s", err.Namespace(), err.Tag(), err.Kind())
				fmt.Println(strings.Replace(msg, "Config.", "", 1))
			}
		} else {
			log.Panic(err)
		}
		os.Exit(1)
	}
}
