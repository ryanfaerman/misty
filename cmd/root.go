package cmd

import (
	"os"

	"github.com/ryanfaerman/misty/config"
	"github.com/ryanfaerman/misty/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logger = log.New()

	configFile string
	verbose    bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is "+config.ConfigFile+")")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logs")
}

var rootCmd = &cobra.Command{
	Use:     config.ApplicationName,
	Short:   config.ApplicationDesc,
	Version: version.Version.String(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			logger.SetLevel(log.DebugLevel)
		}
		config.Logger = logger

		return config.Load(configFile)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
