package cmd

import (
	"github.com/goat-project/goat-one/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var networkRequired = []string{ /* TODO: add required flags here */ }
var networkFlags = []string{ /* TODO: add all network flags here */ }

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Extract network data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about networks, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		checkRequired(networkRequired)
		if viper.GetBool("debug") {
			log.WithFields(log.Fields{"version": version}).Debug("goat-one version")
			logFlags(networkFlags)
		}

		// TODO: do network stuff here
	},
}

func initNetwork() {
	goatOneCmd.AddCommand(networkCmd)

	// TODO: add new flags
	// TODO: configure new flags
}
