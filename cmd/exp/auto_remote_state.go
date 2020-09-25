package exp

import (
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/exp/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	debug bool
	quiet bool
)

func init() {
	autoRemoteStateCmd.Flags().String("path", "", "path to a working directory")
	autoRemoteStateCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable verbose output")
	autoRemoteStateCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "do not output to console; use return code to determine success/failure")

	autoRemoteStateCmd.Flags().StringP("config", "c", "fogg.yml", "Use this to override the fogg config file.")

	ExpCmd.AddCommand(autoRemoteStateCmd)
}

var autoRemoteStateCmd = &cobra.Command{
	Use:   "auto-remote-state",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		setupDebug(debug)
		// Set up fs
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		configFile, e := cmd.Flags().GetString("config")
		if e != nil {
			return errs.WrapInternal(e, "couldn't parse config flag")
		}

		openGitOrExit(fs)

		path, _ := cmd.Flags().GetString("path")
		return state.Run(fs, configFile, path)
	},
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
