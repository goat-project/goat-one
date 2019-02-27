package cmd

import (
	"fmt"
	"strings"

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
		checkRequired(append(vmRequired, append(networkRequired, storageRequired...)...))
		// TODO: do stuff here
	},
}

// Initialize initializes configuration and CLI options.
func Initialize() {
	initGoatOne()
	initVM()
	initNetwork()
	initStorage()
}

func initGoatOne() {
	cobra.OnInitialize(initConfig)

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

	bindFlags(*goatOneCmd, []string{cfgIdentifier, cfgRecordsFrom, cfgRecordsTo, cfgRecordsForPeriod, cfgEndpoint,
		cfgOpennebulaEndpoint, cfgOpennebulaSecret, cfgOpennebulaTimeout, cfgDebug})

	viper.SetDefault("author", "Lenka Svetlovska")
	viper.SetDefault("license", "apache")
}

func initConfig() {
	// name of config file (without extension)
	viper.SetConfigName("goat-one")

	// paths to look for the config file in
	viper.AddConfigPath("config/")
	viper.AddConfigPath("/etc/goat-one/")
	viper.AddConfigPath("$HOME/.goat-one/")

	// find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func checkRequired(required []string) {
	globalRequired := []string{cfgIdentifier, cfgEndpoint, cfgOpennebulaEndpoint, cfgOpennebulaSecret,
		cfgOpennebulaTimeout}

	for _, req := range append(required, globalRequired...) {
		if viper.GetString(req) == "" {
			panic(fmt.Errorf("required flag \"%s\" not set", req))
		}
	}
}

func bindFlags(command cobra.Command, flagsForBinding []string) {
	for _, flag := range flagsForBinding {
		err := viper.BindPFlag(flag, command.PersistentFlags().Lookup(parseFlagName(flag)))
		if err != nil {
			panic(fmt.Errorf("unable to initialize \"%s\" flag", flag))
		}
	}
}

func parseFlagName(CfgName string) string {
	return lastString(strings.Split(CfgName, "."))
}

func lastString(ss []string) string {
	if len(ss) == 0 {
		return ""
	}

	return ss[len(ss)-1]
}
