package integration_test_test

import (
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Startup", func() {
	It("will create the database at startup if it doesn't exist yet", func() {
		RunIntegrationTest(func(exec Executor, cfg *testConfig) {
			err := exec([]string{"--version"})
			Expect(err).ToNot(HaveOccurred())
			Expect(path.Join(cfg.DataDir, "what-next.sqlite")).To(BeAnExistingFile())
		})
	})
})
