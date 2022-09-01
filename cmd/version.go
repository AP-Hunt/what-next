package cmd

import (
	"github.com/spf13/cobra"
)

// Version is set at build time
var Version string = "dev"

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(Version)
	},
}
