package cmd

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cpuprofile string
)

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "do not output to console; use return code to determine success/failure")
	rootCmd.PersistentFlags().StringVarP(&cpuprofile, "cpuprofile", "p", "", "activate cpu profiling via pprof and write to file")
}

var rootCmd = &cobra.Command{
	Use:          "fogg",
	Short:        "",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cpuprofile != "" {
			logrus.Info("starting cpu profile")
			f, err := os.Create(cpuprofile)
			if err != nil {
				return err
			}
			return pprof.StartCPUProfile(f)
		}

		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		quiet, err := cmd.Flags().GetBool("quiet")
		if err != nil {
			return err
		}

		logLevel := logrus.InfoLevel
		if debug { // debug overrides quiet
			logLevel = logrus.DebugLevel
		} else if quiet {
			logLevel = logrus.FatalLevel
		}
		logrus.SetLevel(logLevel)

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if cpuprofile != "" {
			logrus.Info("stopping cpu profile")
			pprof.StopCPUProfile()
		}
		return nil
	},
}

// Execute executes the rootCmd
func Execute() {
	red := color.New(color.FgRed).SprintFunc()

	if err := rootCmd.Execute(); err != nil {
		switch e := err.(type) {
		case *errs.User:
			fmt.Printf("%s: %s\n", red("ERROR"), e.Error())
			os.Exit(1)
		case *errs.Internal:
			fmt.Printf("%s:\nThis may be a bug, please report it.\n\n %s", red("INTERNAL ERROR"), e.Error())
		default:
			fmt.Printf("%s: %s", red("UNKNOWN ERROR"), err)
			os.Exit(1)
		}
	}
}
