package discovery

import (
	"io"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CleanupFunc func()

func TestDiscoverySanity(t *testing.T) {

	fileEntries, cleanup := prepareTestEnvironment()
	defer cleanup()

	bfs, _ := MakeBfsDiscoverer(fileEntries[0].FullPath)

	for i := 0; i < len(fileEntries); i++ {
		res, err := bfs.Next()
		assert.Nil(t, err)
		assert.Equal(t, fileEntries[i], res)
	}

	res, err := bfs.Next()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, FileInformation{}, res)
}

func prepareTestEnvironment() ([]FileInformation, CleanupFunc) {
	// Create a well known directory structure and test that all entries are correctly detected
	baseDir, _ := os.MkdirTemp("", "")
	myFile := path.Join(baseDir, "myfile.txt")
	myFile2 := path.Join(baseDir, "myfile2.txt")
	backupDir := path.Join(baseDir, "backup")
	backupFile := path.Join(backupDir, "bak")
	os.WriteFile(myFile, []byte{1, 2, 3, 4, 5, 6}, 0)
	os.WriteFile(myFile2, []byte{1, 2, 3, 4, 5, 6}, 0)
	_ = os.Mkdir(backupDir, fs.ModePerm)
	os.WriteFile(backupFile, []byte{1, 2, 3, 4, 5, 6}, 0)

	baseDirInfo, _ := MakeFileInformation(baseDir)
	myFileInfo, _ := MakeFileInformation(myFile)
	myFile2Info, _ := MakeFileInformation(myFile2)
	backupDirInfo, _ := MakeFileInformation(backupDir)
	backupFileInfo, _ := MakeFileInformation(backupFile)

	cleanup := func() {
		os.RemoveAll(baseDir)
	}

	return []FileInformation{
		baseDirInfo,
		myFileInfo,
		myFile2Info,
		backupDirInfo,
		backupFileInfo,
	}, cleanup
}
