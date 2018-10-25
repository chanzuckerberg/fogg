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
	validator "gopkg.in/go-playground/validator.v9"
)

func openGitOrExit(pwd string) {
	log.Debugf("opening git at %s", pwd)
	// g, err := git.PlainOpen(pwd)
	_, err := os.Stat(".git")
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
