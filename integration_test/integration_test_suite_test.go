package integration_test_test

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdContext "github.com/AP-Hunt/what-next/m/context"
)

var binaryLocation string

func init() {
	flag.StringVar(&binaryLocation, "binary", "./bin/what-next", "Specify the path of the binary to test")
}

type testConfig struct {
	Context cmdContext.CommandContext
	DataDir string
}

type Executor = func(args []string) error

func RunIntegrationTest(action func(exec Executor, cfg *testConfig)) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "what-next-integration-test")
	Expect(err).ToNot(HaveOccurred())

	defer os.RemoveAll(tempDir)

	os.Setenv("WHAT_NEXT_DATA_DIR", ":memory:")

	Expect(path.Join(tempDir, "what-next.sqlite")).
		ToNot(BeAnExistingFile(), "The database file should not exist at the start of integration testing")

	cfg := &testConfig{
		DataDir: tempDir,
	}

	executor := func(args []string) error {
		binPath, err := filepath.Abs(binaryLocation)
		Expect(err).ToNot(HaveOccurred())

		cmd := exec.Command(binPath, args...)

		env := os.Environ()
		env = append(env, fmt.Sprintf("WHAT_NEXT_DATA_DIR=%s", tempDir))

		errBuff := bytes.Buffer{}
		errBuffWriter := bufio.NewWriter(&errBuff)

		cmd.Env = env
		cmd.Stdout = GinkgoWriter
		cmd.Stderr = errBuffWriter

		runErr := cmd.Run()

		if runErr != nil {
			stdErrContent := errBuff.String()
			return fmt.Errorf("%s\n%s\n", runErr, stdErrContent)
		}

		return nil
	}

	action(executor, cfg)
}

func TestIntegrationTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTest Suite")
}
