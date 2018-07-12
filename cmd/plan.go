package cmd

import (
	"fmt"
	"os"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	planCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	planCmd.Flags().BoolP("sicc", "s", false, "Use this to turn on sicc-compatibility mode. Implies -c sicc.json.")
	planCmd.Flags().BoolP("verbose", "v", false, "use this to turn on verbose output")
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run a plan",
	Run: func(cmd *cobra.Command, args []string) {
		var e error
		// Set up fs

		pwd, e := os.Getwd()
		if e != nil {
			panic(e)
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		siccMode, e := cmd.Flags().GetBool("sicc")
		if e != nil {
			panic(e)
		}
		verbose, e := cmd.Flags().GetBool("verbose")
		if e != nil {
			panic(e)
		}
		var configFile string
		if siccMode {
			configFile = "sicc.json"
		} else {
			configFile, e = cmd.Flags().GetString("config")
			if e != nil {
				panic(e)
			}
		}

		// check that we are at root of initialized git repo
		openGitOrExit(pwd)

		config, err := config.FindAndReadConfig(fs, configFile)
		if err != nil {
			panic(err)
		}
		if verbose {
			fmt.Println("CONFIG")
			fmt.Printf("%#v\n=====", config)
		}

		err = config.Validate()

		if err != nil {
			fmt.Println("Config Error(s):")
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Printf("\t%s is a %s %s\n", err.Namespace(), err.Tag(), err.Kind())
			}

			os.Exit(1)
		}

		p, e := plan.Eval(config, siccMode, verbose)
		if e != nil {
			panic(e)
		}
		e = plan.Print(p)
		if e != nil {
			panic(e)
		}
	},
}
