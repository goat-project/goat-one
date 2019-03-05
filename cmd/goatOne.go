package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/goat-project/goat-one/constants"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "1.0.0"

var goatOneCmd = &cobra.Command{
	Use:   "goat-one",
	Short: "extracts data about virtual machines, networks and storages",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about virtual machines, networks and storages, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		checkRequired(append(vmRequired, append(networkRequired, storageRequired...)...))
		// TODO: do stuff here
	},
}

// Execute uses the args (os.Args[1:] by default)
// and run through the command tree finding appropriate matches
// for commands and then corresponding flags.
func Execute() {
	if err := goatOneCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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

	goatOneCmd.PersistentFlags().StringP(constants.CfgIdentifier, "i", viper.GetString(constants.CfgIdentifier),
		"goat identifier [IDENTIFIER] (required)")
	goatOneCmd.PersistentFlags().StringP(constants.CfgRecordsFrom, "f", viper.GetString(constants.CfgRecordsFrom),
		"records from [TIME]")
	goatOneCmd.PersistentFlags().StringP(constants.CfgRecordsTo, "t", viper.GetString(constants.CfgRecordsTo),
		"records to [TIME]")
	goatOneCmd.PersistentFlags().StringP(constants.CfgRecordsForPeriod, "p",
		viper.GetString(constants.CfgRecordsForPeriod), "records for period [TIME PERIOD]")
	goatOneCmd.PersistentFlags().StringP(constants.CfgEndpoint, "e", viper.GetString(constants.CfgEndpoint),
		"goat server [GOAT_SERVER_ENDPOINT] (required)")
	goatOneCmd.PersistentFlags().StringP(constants.CfgOpennebulaEndpoint, "o",
		viper.GetString(constants.CfgOpennebulaEndpoint), "OpenNebula endpoint [OPENNEBULA_ENDPOINT] (required)")
	goatOneCmd.PersistentFlags().StringP(constants.CfgOpennebulaSecret, "s",
		viper.GetString(constants.CfgOpennebulaSecret), "OpenNebula secret [OPENNEBULA_SECRET] (required)")
	goatOneCmd.PersistentFlags().String(constants.CfgOpennebulaTimeout, viper.GetString(constants.CfgOpennebulaTimeout),
		"timeout for OpenNebula calls [TIMEOUT_FOR_OPENNEBULA_CALLS] (required)")
	goatOneCmd.PersistentFlags().StringP(constants.CfgDebug, "d", viper.GetString(constants.CfgDebug),
		"debug")

	bindFlags(*goatOneCmd, []string{constants.CfgIdentifier, constants.CfgRecordsFrom, constants.CfgRecordsTo,
		constants.CfgRecordsForPeriod, constants.CfgEndpoint, constants.CfgOpennebulaEndpoint,
		constants.CfgOpennebulaSecret, constants.CfgOpennebulaTimeout, constants.CfgDebug})

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
		// TODO log that configuration file couldn't be read
		fmt.Printf("error config file: %s", err)
	}
}

func checkRequired(required []string) {
	globalRequired := []string{constants.CfgIdentifier, constants.CfgEndpoint, constants.CfgOpennebulaEndpoint,
		constants.CfgOpennebulaSecret, constants.CfgOpennebulaTimeout}

	for _, req := range append(required, globalRequired...) {
		if viper.GetString(req) == "" {
			// TODO log that required flag is missing
			fmt.Printf("required flag \"%s\" not set", req)
			os.Exit(1)
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

func parseFlagName(cfgName string) string {
	return lastString(strings.Split(cfgName, "."))
}

func lastString(ss []string) string {
	// This should not happen since it is passing a predefined non-empty strings.
	// It panic here since this will happen only if a mistake in code is made.
	if len(ss) == 0 {
		panic("parsing empty string")
	}

	return ss[len(ss)-1]
}
