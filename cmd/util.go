package cmd

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"github.com/chanzuckerberg/fogg/config"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	validator "gopkg.in/go-playground/validator.v9"
)

func openGitOrExit(fs afero.Fs) {
	_, err := fs.Stat(".git")
	if err != nil {
		// assuming this means no repository
		logrus.Fatal("fogg must be run from the root of a git repo")
		os.Exit(1)
	}
}

func readAndValidateConfig(fs afero.Fs, configFile string) (*v2.Config, []string, error) {
	conf, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return nil, nil, errs.WrapUser(err, "unable to read config file")
	}
	logrus.Debug("CONFIG")
	logrus.Debugf("%s\n=====", pretty.Sprint(conf))

	warnings, e := conf.Validate()
	return conf, warnings, e
}

func mergeConfigValidationErrors(err error) error {
	if err != nil {
		fmt.Println("fogg.yml has error(s):")
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
	logLevel := logrus.InfoLevel
	if debug { // debug overrides quiet
		logLevel = logrus.DebugLevel
		go func() {
			logrus.Println(http.ListenAndServe("localhost:6060", nil))
			http.HandleFunc("/", pprof.Index)
		}()
	} else if quiet {
		logLevel = logrus.FatalLevel
	}
	logrus.SetLevel(logLevel)
}

func openFs() (afero.Fs, error) {
	pwd, e := os.Getwd()
	if e != nil {
		return nil, e
	}
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	return fs, nil
}

func bootstrapCmd(cmd *cobra.Command, debug bool) (afero.Fs, *v2.Config, error) {
	setupDebug(debug)

	fs, err := openFs()
	if err != nil {
		return nil, nil, err
	}

	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, nil, err
	}

	config, _, err := readAndValidateConfig(fs, configFile)

	err = mergeConfigValidationErrors(err)
	if err != nil {
		return nil, nil, err
	}

	return fs, config, nil
}

func printWarnings(warnings []string) {
	for _, w := range warnings {
		logrus.Warn(w)
	}
}
