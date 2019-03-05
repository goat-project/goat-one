package cmd

import (
	"strings"

	"github.com/goat-project/goat-one/logger"

	"github.com/goat-project/goat-one/constants"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const version = "1.0.0"

var goatOneFlags = []string{constants.CfgIdentifier, constants.CfgRecordsFrom, constants.CfgRecordsTo,
	constants.CfgRecordsForPeriod, constants.CfgEndpoint, constants.CfgOpennebulaEndpoint,
	constants.CfgOpennebulaSecret, constants.CfgOpennebulaTimeout, constants.CfgDebug, constants.CfgLogPath}

var goatOneCmd = &cobra.Command{
	Use:   "goat-one",
	Short: "extracts data about virtual machines, networks and storages",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about virtual machines, networks and storages, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		checkRequired(append(vmRequired, append(networkRequired, storageRequired...)...))
		if viper.GetBool("debug") {
			log.WithFields(log.Fields{"version": version}).Debug("goat-one version")
			logFlags(append(vmFlags, append(networkFlags, storageFlags...)...))
		}
		// TODO: do stuff here
	},
}

// Execute uses the args (os.Args[1:] by default)
// and run through the command tree finding appropriate matches
// for commands and then corresponding flags.
func Execute() {
	if err := goatOneCmd.Execute(); err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("fatal error execute")
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
	goatOneCmd.PersistentFlags().String(constants.CfgLogPath, viper.GetString(constants.CfgLogPath), "path to log file")

	bindFlags(*goatOneCmd, goatOneFlags)

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
		log.WithFields(log.Fields{"error": err}).Error("error config file")
	}
}

func checkRequired(required []string) {
	globalRequired := []string{constants.CfgIdentifier, constants.CfgEndpoint, constants.CfgOpennebulaEndpoint,
		constants.CfgOpennebulaSecret, constants.CfgOpennebulaTimeout}

	for _, req := range append(required, globalRequired...) {
		if viper.GetString(req) == "" {
			log.WithFields(log.Fields{"flag": req}).Fatal("required flag not set")
		}
	}
}

func bindFlags(command cobra.Command, flagsForBinding []string) {
	for _, flag := range flagsForBinding {
		err := viper.BindPFlag(flag, command.PersistentFlags().Lookup(parseFlagName(flag)))
		if err != nil {
			log.WithFields(log.Fields{"error": err, "flag": flag}).Panic("unable to initialize flag")
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
		log.Panic("parsing empty string")
	}

	return ss[len(ss)-1]
}

func logFlags(flags []string) {
	for _, flag := range append(goatOneFlags, flags...) {
		log.WithFields(log.Fields{"flag": flag, "value": viper.Get(flag)}).Debug("flag initialized")
	}
}
