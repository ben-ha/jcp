package copier

import (
	"io"
	"os"
	"os/exec"
	"testing"

	discovery "github.com/ben-ha/jcp/discovery"
	"github.com/stretchr/testify/assert"
)

func TestCopyBusyFile(t *testing.T) {
	lsPath := "/bin/ls"
	runningProcessPath := copyFileToTemp("/bin/cat")
	os.Chmod(runningProcessPath, 0700)
	runningProcess := exec.Command(runningProcessPath)
	runningProcess.Start()
	defer runningProcess.Process.Kill()

	copier := BlockCopier{BlockSize: 512}
	sourceInfo, _ := discovery.MakeFileInformation(lsPath)
	destInfo, _ := discovery.MakeFileInformation(runningProcessPath)
	newState := copier.Copy(sourceInfo, destInfo, CopierState{State: BlockCopierState{}})

	assert.NotNil(t, newState.Error)
	assert.Equal(t, io.EOF.Error(), newState.Error.Error())
}

func copyFileToTemp(path string) string {
	f, _ := os.CreateTemp("", "block_copier_test")
	defer f.Close()
	src, _ := os.Open(path)
	defer src.Close()
	_, _ = io.Copy(f, src)
	return f.Name()
}
