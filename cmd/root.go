package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/chanzuckerberg/fogg/cmd/exp"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	debug      bool
	quiet      bool
	cpuprofile string
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "do not output to console; use return code to determine success/failure")
	rootCmd.PersistentFlags().StringVarP(&cpuprofile, "cpuprofile", "p", "", "activate cpu profiling via pprof and write to file")
	rootCmd.AddCommand(exp.ExpCmd)
}

var rootCmd = &cobra.Command{
	Use:          "fogg",
	Short:        "",
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
		if cpuprofile != "" {
			log.Println("starting cpu profile")
			f, err := os.Create(cpuprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
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
			fmt.Printf("%s: %s", red("UNKOWN ERROR"), err)
			os.Exit(1)
		}
	}
}
