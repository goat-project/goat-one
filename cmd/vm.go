package cmd

import (
	"time"

	"github.com/goat-project/goat-one/filter"

	"github.com/goat-project/goat-one/client"
	"github.com/goat-project/goat-one/constants"
	"github.com/goat-project/goat-one/logger"
	"github.com/goat-project/goat-one/preparer"
	"github.com/goat-project/goat-one/processor"
	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/resource/virtualmachine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"

	log "github.com/sirupsen/logrus"
)

var vmRequired = []string{constants.CfgSiteName, constants.CfgCloudType}
var vmFlags = []string{constants.CfgSiteName, constants.CfgCloudType, constants.CfgCloudComputeService}

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Extract virtual machine data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about virtual machines, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		checkRequired(vmRequired)
		if viper.GetBool("debug") {
			log.WithFields(log.Fields{"version": version}).Debug("goat-one version")
			logFlags(vmFlags)
		}

		readLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)
		writeLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond)

		accountVM(readLimiter, writeLimiter)
	},
}

func initVM() {
	goatOneCmd.AddCommand(vmCmd)

	vmCmd.PersistentFlags().String(parseFlagName(constants.CfgSiteName), viper.GetString(constants.CfgSiteName),
		"site name [VM_SITE_NAME] (required)")
	vmCmd.PersistentFlags().String(parseFlagName(constants.CfgCloudType), viper.GetString(constants.CfgCloudType),
		"cloud type [VM_CLOUD_TYPE] (required)")
	vmCmd.PersistentFlags().String(parseFlagName(constants.CfgCloudComputeService),
		viper.GetString(constants.CfgCloudComputeService), "cloud compute service [VM_CLOUD_COMPUTE_SERVICE]")

	bindFlags(*vmCmd, vmFlags)
}

func accountVM(readLimiter, writeLimiter *rate.Limiter) {
	read := reader.CreateReader(getOpenNebulaClient(), readLimiter)

	proc := processor.CreateProcessor(virtualmachine.CreateProcessor(read))
	filt := filter.CreateFilter(virtualmachine.CreateFilter())
	prep := preparer.CreatePreparer(virtualmachine.CreatePreparer(read, writeLimiter, getConn()))

	c := client.Client{}

	c.Run(proc, filt, prep)
}
