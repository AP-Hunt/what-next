package cmd

import (
	. "github.com/AP-Hunt/what-next/m/context"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "what-next",
	Version: Version,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(todoCmd)
}

func ExecuteCWithArgs(ctx CommandContext, args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.ExecuteContext(ctx)
}

func ExecuteC(ctx CommandContext) error {
	return rootCmd.ExecuteContext(ctx)
}
