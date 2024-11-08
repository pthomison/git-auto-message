package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "unset" // default value
var commit = "unset"  // default value
var date = "unset"    // default value

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %v\n", version)
		fmt.Printf("Commit: %v\n", commit)
		fmt.Printf("Date: %v\n", date)
	},
}
