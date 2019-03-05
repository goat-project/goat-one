package cmd

import (
	"github.com/spf13/cobra"
)

var networkRequired = []string{ /* TODO: add required flags here */ }

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Extract network data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about networks, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		checkRequired(networkRequired)
		// TODO: do network stuff here
	},
}

func initNetwork() {
	goatOneCmd.AddCommand(networkCmd)

	// TODO: add new flags
	// TODO: configure new flags
}
