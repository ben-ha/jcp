package copier

import (
	"io"
	"os"
	"testing"

	discovery "github.com/ben-ha/jcp/discovery"
	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	expectedData := "Hello world!"
	sourceFile := prepareFile(expectedData)
	destFile := prepareTemporaryFileName()

	copier := BlockCopier{BlockSize: 512}
	sourceInfo, _ := discovery.MakeFileInformation(sourceFile)
	newState := copier.Copy(sourceInfo, discovery.FileInformation{FullPath: destFile, Info: nil}, CopierState{State: BlockCopierState{}})

	assert.NotNil(t, newState.Error)
	assert.Equal(t, newState.Error, io.EOF)

	actualDataBytes, _ := os.ReadFile(destFile)
	actualData := string(actualDataBytes)

	if actualData != expectedData {
		t.Fatalf("Received different data. Expected=%v, Actual=%v", expectedData, actualData)
	}
}

func TestCopyFileProgress(t *testing.T) {
	expectedData := "Hello world!"
	sourceFile := prepareFile(expectedData)
	destFile := prepareTemporaryFileName()

	copier := BlockCopier{BlockSize: 1}
	progressChannel := make(chan CopierProgress, 100)
	sourceInfo, _ := discovery.MakeFileInformation(sourceFile)
	newState := copier.CopyWithProgress(sourceInfo, discovery.FileInformation{FullPath: destFile, Info: nil}, CopierState{State: BlockCopierState{}}, progressChannel)

	assert.NotNil(t, newState.Error)
	assert.Equal(t, newState.Error, io.EOF)
	close(progressChannel)

	currentBlockTransferred := uint64(0)
	for data := range progressChannel {
		assert.Equal(t, data.BytesTransferred, data.BytesTransferred)
		assert.Equal(t, uint64(len(expectedData)), data.Size)
		assert.Equal(t, sourceFile, data.Source)
		assert.Equal(t, destFile, data.Dest)
		currentBlockTransferred++
	}
}

func TestResumeCopy(t *testing.T) {
	expectedData := "Hello world!"
	partialData := expectedData[0:3]
	sourceFile := prepareFile(expectedData)
	partialDestFile := prepareFile(partialData)

	copier := BlockCopier{BlockSize: 1}
	copierState := CopierState{State: BlockCopierState{Size: uint64(len(expectedData)), BytesTransferred: uint64(len(partialData))}}

	sourceInfo, _ := discovery.MakeFileInformation(sourceFile)
	newState := copier.Copy(sourceInfo, discovery.FileInformation{FullPath: partialDestFile, Info: MakeFakeDestinationFileInfo(partialDestFile, 3)}, copierState)

	assert.NotNil(t, newState.Error)
	assert.Equal(t, newState.Error, io.EOF)

	actualDataBytes, _ := os.ReadFile(partialDestFile)
	actualData := string(actualDataBytes)

	if actualData != expectedData {
		t.Fatalf("Received different data. Expected=%v, Actual=%v", expectedData, actualData)
	}
}

func TestResumeDeletedFileCopy(t *testing.T) {
	expectedData := "Hello world!"
	partialData := expectedData[0:3]
	sourceFile := prepareFile(expectedData)
	destFile := prepareTemporaryFileName()

	copier := BlockCopier{BlockSize: 1}
	copierState := CopierState{State: BlockCopierState{Size: uint64(len(expectedData)), BytesTransferred: uint64(len(partialData))}}

	sourceInfo, _ := discovery.MakeFileInformation(sourceFile)
	newState := copier.Copy(sourceInfo, discovery.FileInformation{FullPath: destFile, Info: MakeFakeDestinationFileInfo(destFile, 0)}, copierState)

	assert.NotNil(t, newState.Error)
	assert.Equal(t, newState.Error, io.EOF)

	actualDataBytes, _ := os.ReadFile(destFile)
	actualData := string(actualDataBytes)

	if actualData != expectedData {
		t.Fatalf("Received different data. Expected=%v, Actual=%v", expectedData, actualData)
	}
}

func prepareFile(data string) string {
	f, _ := os.CreateTemp("", "block_copier_test")
	_, _ = f.WriteString(data)
	_ = f.Close()
	return f.Name()
}

func prepareTemporaryFileName() string {
	f, _ := os.CreateTemp("", "block_copier_test")
	_ = f.Close()
	os.Remove(f.Name())
	return f.Name()
}
