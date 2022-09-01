package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:     "what-next",
	Version: Version,
}

func Execute() error {
	rootCmd.AddCommand(versionCmd)

	return rootCmd.Execute()
}
