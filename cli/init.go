package cli

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	global *bool
)

func init()  {
	rootCmd.AddCommand(initCmd)
	global = initCmd.Flags().BoolP("global", "g", false, "Should the repo be initialized globally")
}

var initCmd = &cobra.Command{
	Use: "init",
	Short: "Initializes the node and generates the config file",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(runInit(args))
	},
}

func runInit(args []string) int {
	var home string

	if *global {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return printError(err)
		}
	} else {
		home = "."
	}

	// initialize the config
	initializeConfig(home)
	// initialize the store directory
	storeDirectory(home)

	return 0

}
