package discovery

import (
	"io/fs"
	"os"
)

type FileInformation struct {
	FullPath string
	Info fs.FileInfo
}

func MakeFileInformation(fullPath string) (FileInformation, error) {
	stat, statErr := os.Stat(fullPath)
	if statErr != nil {
		return FileInformation{}, statErr
	}

	return FileInformation{FullPath: fullPath, Info:stat}, nil
}