package cmd

import "github.com/spf13/cobra"

const (
// TODO: add constants for flags here
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Extract storage data",
	Long: "The accounting client is a command-line tool that connects to a cloud, " +
		"extracts data about storages, filters them accordingly and " +
		"then sends them to a server for further processing.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: do storage stuff here
	},
}

func initStorage() {
	goatOneCmd.AddCommand(storageCmd)

	// TODO: add new flags
	// TODO: configure new flags
}
