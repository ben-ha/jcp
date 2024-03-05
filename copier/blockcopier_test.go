package copier

import (
	"os"
	"testing"
	discovery "github.com/ben-ha/jcp/discovery"
)

func TestCopyFile(t *testing.T) {
	expectedData := "Hello world!"
	sourceFile := prepareFile(expectedData)
	destFile := prepareTemporaryFileName()

	copier := BlockCopier{BlockSize: 512}
	sourceInfo, _ := discovery.MakeFileInformation(sourceFile)
	newState := copier.Copy(sourceInfo, discovery.FileInformation{FullPath: destFile, Info: nil}, CopierState{State: BlockCopierState{}})

	if newState.Error != nil {
		t.Fatalf("Unexpected error occurred: %v", (*newState.Error).Error())
	}

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

	if newState.Error != nil {
		t.Fatalf("Unexpected error occurred: %v", (*newState.Error).Error())
	}

	currentBlockTransferred := uint64(0)
	for data := range progressChannel {
		if data.BytesTransferred != currentBlockTransferred {
			t.Fatalf("Unexpected progress. Expected=%v, actual=%v", currentBlockTransferred, data.BytesTransferred)
		}
		if data.Size != uint64(len(expectedData)) {
			t.Fatalf("Unexpected size. Expected = %v, actual = %v", len(expectedData), data.Size)
		}
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

	if newState.Error != nil {
		t.Fatalf("Unexpected error occurred: %v", (*newState.Error).Error())
	}

	actualDataBytes, _ := os.ReadFile(partialDestFile)
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
