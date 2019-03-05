package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cfgVMSubcommand = "vm."

const (
	cfgSiteName            = cfgVMSubcommand + "site-name"
	cfgCloudType           = cfgVMSubcommand + "cloud-type"
	cfgCloudComputeService = cfgVMSubcommand + "cloud-compute-service"
)

var vmRequired = []string{cfgSiteName, cfgCloudType}

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Extract virtual machine data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about virtual machines, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		checkRequired(vmRequired)
		// TODO: do VM stuff here
	},
}

func initVM() {
	goatOneCmd.AddCommand(vmCmd)

	vmCmd.PersistentFlags().String(parseFlagName(cfgSiteName), viper.GetString(cfgSiteName),
		"site name [VM_SITE_NAME] (required)")
	vmCmd.PersistentFlags().String(parseFlagName(cfgCloudType), viper.GetString(cfgCloudType),
		"cloud type [VM_CLOUD_TYPE] (required)")
	vmCmd.PersistentFlags().String(parseFlagName(cfgCloudComputeService), viper.GetString(cfgCloudComputeService),
		"cloud compute service [VM_CLOUD_COMPUTE_SERVICE]")

	bindFlags(*vmCmd, []string{cfgSiteName, cfgCloudType, cfgCloudComputeService})
}
