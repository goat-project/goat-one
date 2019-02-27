package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgIdentifier         = "identifier"
	cfgRecordsFrom        = "records-from"
	cfgRecordsTo          = "records-to"
	cfgRecordsForPeriod   = "records-for-period"
	cfgEndpoint           = "endpoint"
	cfgOpennebulaEndpoint = "opennebula-endpoint"
	cfgOpennebulaSecret   = "opennebula-secret" // nolint: gosec
	cfgOpennebulaTimeout  = "opennebula-timeout"
	cfgDebug              = "debug"
	cfgVersion            = "version"
)

var goatOneCmd = &cobra.Command{
	Use:   "goat-one",
	Short: "extracts data about virtual machines, networks and storages",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about virtual machines, networks and storages, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Version: viper.GetString(cfgVersion),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: do stuff here
	},
}

func initGoatOne() {
	goatOneCmd.PersistentFlags().StringP(cfgIdentifier, "i", viper.GetString(cfgIdentifier),
		"goat identifier [IDENTIFIER] (required)")
	goatOneCmd.PersistentFlags().StringP(cfgRecordsFrom, "f", viper.GetString(cfgRecordsFrom),
		"records from [TIME]")
	goatOneCmd.PersistentFlags().StringP(cfgRecordsTo, "t", viper.GetString(cfgRecordsTo),
		"records to [TIME]")
	goatOneCmd.PersistentFlags().StringP(cfgRecordsForPeriod, "p", viper.GetString(cfgRecordsForPeriod),
		"records for period [TIME PERIOD]")
	goatOneCmd.PersistentFlags().StringP(cfgEndpoint, "e", viper.GetString(cfgEndpoint),
		"goat server [GOAT_SERVER_ENDPOINT] (required)")
	goatOneCmd.PersistentFlags().StringP(cfgOpennebulaEndpoint, "o", viper.GetString(cfgOpennebulaEndpoint),
		"OpenNebula endpoint [OPENNEBULA_ENDPOINT] (required)")
	goatOneCmd.PersistentFlags().StringP(cfgOpennebulaSecret, "s", viper.GetString(cfgOpennebulaSecret),
		"OpenNebula secret [OPENNEBULA_SECRET] (required)")
	goatOneCmd.PersistentFlags().String(cfgOpennebulaTimeout, viper.GetString(cfgOpennebulaTimeout),
		"timeout for OpenNebula calls [TIMEOUT_FOR_OPENNEBULA_CALLS] (required)")
	goatOneCmd.PersistentFlags().StringP(cfgDebug, "d", viper.GetString(cfgDebug),
		"debug")
}
