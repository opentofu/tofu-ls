package exec

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opentofu/tofudl"
)

func NewTestingExecutor(t *testing.T, workDir string) TerraformExecutor {

	dl, err := tofudl.New()
	if err != nil {
		log.Fatalf("error when instantiating tofudl %s", err)
	}

	binary, err := dl.Download(t.Context())
	if err != nil {
		log.Fatalf("error when downloading %s", err)
	}

	execPath := filepath.Join(workDir, "tofu")
	// Windows executable case
	if runtime.GOOS == "windows" {
		execPath += ".exe"
	}
	if err := os.WriteFile(execPath, binary, 0755); err != nil {
		log.Fatalf("error when writing the file %s: %s", execPath, err)
	}

	t.Cleanup(func() {
		if err := os.Remove(execPath); err != nil {
			t.Fatal(err)
		}
	})

	e, err := NewExecutor(workDir, execPath)
	if err != nil {
		t.Fatal(err)
	}
	return e
}
