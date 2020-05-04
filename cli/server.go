package cli

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use: "server [command]",
	Short: "Interact with your local server",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
