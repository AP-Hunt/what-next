package cmd

import (
	. "github.com/AP-Hunt/what-next/m/context"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "what-next",
	Version: Version,
}

func init() {
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(TodoRootCmd)
}

func ExecuteCWithArgs(ctx CommandContext, args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.ExecuteContext(ctx)
}

func ExecuteC(ctx CommandContext) error {
	return RootCmd.ExecuteContext(ctx)
}
