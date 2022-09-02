package integration_test_test

import (
	"context"
	"os"
	"path"

	"github.com/AP-Hunt/what-next/m/cmd"
	cmdContext "github.com/AP-Hunt/what-next/m/context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Startup", func() {
	It("will create the database at startup if it doesn't exist yet", func() {

		tempDir, err := os.MkdirTemp(os.TempDir(), "what-next-integration-test")
		Expect(err).ToNot(HaveOccurred())
		defer os.RemoveAll(tempDir)

		os.Setenv("WHAT_NEXT_DATA_DIR", tempDir)

		Expect(path.Join(tempDir, "what-next.sqlite")).ToNot(BeAnExistingFile())

		ctx, err := cmdContext.CreateDefaultCommandContext(context.Background())
		Expect(err).ToNot(HaveOccurred())
		cmd.ExecuteCWithArgs(ctx, []string{"--version"})

		Expect(path.Join(tempDir, "what-next.sqlite")).To(BeAnExistingFile())
	})
})
