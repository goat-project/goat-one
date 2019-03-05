package cmd

import (
	"github.com/goat-project/goat-one/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var storageRequired = []string{ /* TODO: add required flags here */ }
var storageFlags = []string{ /* TODO: add all storage flags here */ }

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Extract storage data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about storages, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		checkRequired(storageRequired)
		if viper.GetBool("debug") {
			log.WithFields(log.Fields{"version": version}).Debug("goat-one version")
			logFlags(storageFlags)
		}

		// TODO: do storage stuff here
	},
}

func initStorage() {
	goatOneCmd.AddCommand(storageCmd)

	// TODO: add new flags
	// TODO: configure new flags
}
