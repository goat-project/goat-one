package cmd

import (
	"time"

	"github.com/goat-project/goat-one/filter"
	"github.com/goat-project/goat-one/preparer"
	"github.com/goat-project/goat-one/processor"

	"github.com/goat-project/goat-one/client"
	"github.com/goat-project/goat-one/logger"
	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/resource/network"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"

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

		readLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)
		writeLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)

		accountNetwork(readLimiter, writeLimiter)
	},
}

func initNetwork() {
	goatOneCmd.AddCommand(networkCmd)

	// TODO: add new flags
	// TODO: configure new flags
}

func accountNetwork(readLimiter, writeLimiter *rate.Limiter) {
	read := reader.CreateReader(readLimiter)

	prep := preparer.CreatePreparer(network.CreatePreparer(writeLimiter))
	filt := filter.CreateFilter(network.CreateFilter())
	proc := processor.CreateProcessor(network.CreateProcessor(read))

	c := client.Client{}

	c.Run(proc, filt, prep)
}
