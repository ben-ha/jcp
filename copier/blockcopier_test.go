package copier

import (
	"testing"
	"os"
)

func TestCopyFile(t *testing.T) {
	expectedData := "Hello world!"
	sourceFile := prepareFile(expectedData)
	destFile := prepareTemporaryFileName()

	copier := BlockCopier{ BlockSize : 512 }
	newState := copier.Copy(sourceFile, destFile, CopierState{State: BlockCopierState{}})

	if newState.Error != nil {
		t.Fatalf("Unexpected error occurred: %v", (*newState.Error).Error())
	}
	
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
