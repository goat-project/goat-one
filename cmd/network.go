package cmd

import (
	"github.com/spf13/cobra"
)

const (
// TODO: add constants for flags here
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Extract network data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about networks, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: do network stuff here
	},
}

func initNetwork() {
	goatOneCmd.AddCommand(networkCmd)

	// TODO: add new flags
	// TODO: configure new flags
}
