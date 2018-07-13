package cmd

import (
	"fmt"
	"os"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/go-playground/validator"
	"github.com/spf13/afero"
	git "gopkg.in/src-d/go-git.v4"
)

func openGitOrExit(pwd string) *git.Repository {
	g, err := git.PlainOpen(pwd)
	if err != nil {
		// assuming this means no repository
		fmt.Println("fogg must be run from the root of a git repo")
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
		fmt.Println("CONFIG")
		fmt.Printf("%#v\n=====", config)
	}

	err = config.Validate()
	return config, err
}

func exitOnConfigErrors(err error) {
	if err != nil {
		fmt.Println("Config Error(s):")
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				fmt.Printf("\t%s is a %s %s\n", err.Namespace(), err.Tag(), err.Kind())
			}
		} else {
			panic(err)
		}
		os.Exit(1)
	}
}
