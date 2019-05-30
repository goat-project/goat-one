package cmd

import (
	"time"

	"github.com/goat-project/goat-one/constants"

	"github.com/goat-project/goat-one/client"
	"github.com/goat-project/goat-one/filter"
	"github.com/goat-project/goat-one/logger"
	"github.com/goat-project/goat-one/preparer"
	"github.com/goat-project/goat-one/processor"
	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/resource/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"

	log "github.com/sirupsen/logrus"
)

var storageRequired = []string{}
var storageFlags = []string{constants.CfgSite}

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

		readLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)
		writeLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)

		accountStorage(readLimiter, writeLimiter)
	},
}

func initStorage() {
	goatOneCmd.AddCommand(storageCmd)

	storageCmd.PersistentFlags().String(parseFlagName(constants.CfgSite),
		viper.GetString(constants.CfgSite), "site [SITE]")

	bindFlags(*storageCmd, storageFlags)
}

func accountStorage(readLimiter, writeLimiter *rate.Limiter) {
	read := reader.CreateReader(readLimiter)

	proc := processor.CreateProcessor(storage.CreateProcessor(read))
	filt := filter.CreateFilter(storage.CreateFilter())
	prep := preparer.CreatePreparer(storage.CreatePreparer(writeLimiter))

	c := client.Client{}

	c.Run(proc, filt, prep)
}
