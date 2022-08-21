package file

import (
	_ "embed"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"io"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

//go:embed test-fixtures/template.pkr.hcl
var testDatasourceBasic string

func TestAccSopsFile(t *testing.T) {
	const expectedFileContent = "secret: Hello, World"

	testCase := &acctest.PluginTestCase{
		Name: "sops_file_datasource_basic_test",
		Setup: func() error {
			return os.Setenv("SOPS_AGE_KEY_FILE", "test-fixtures/test_key.txt")
		},
		Teardown: func() error {
			return nil
		},
		Template: testDatasourceBasic,
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := io.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			logsString := string(logsBytes)

			valueLog := fmt.Sprintf("null.basic-example: decrypted value: %s", expectedFileContent)

			if matched, _ := regexp.MatchString(valueLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected arn %q", logsString)
			}

			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
