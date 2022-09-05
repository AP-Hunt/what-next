package cmd_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

func PrepareCommandForTest(command *cobra.Command, args []string) {
	command.ResetCommands()
	command.ResetFlags()
	command.SetArgs(args)
	command.SetOutput(GinkgoWriter)
}
