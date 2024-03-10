package discovery

import (
	"io/fs"
	"os"
	"path/filepath"
)

type FileInformation struct {
	FullPath string
	Info     fs.FileInfo
}

func MakeFileInformation(path string) (FileInformation, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return FileInformation{FullPath: path, Info: nil}, err
	}
	stat, statErr := os.Stat(absPath)
	if statErr != nil {
		return FileInformation{FullPath: absPath, Info: nil}, statErr
	}

	return FileInformation{FullPath: absPath, Info: stat}, nil
}

func MakeFileInformationWithSymbolicLinks(path string) (FileInformation, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return FileInformation{FullPath: path, Info: nil}, err
	}
	stat, statErr := os.Lstat(absPath)
	if statErr != nil {
		return FileInformation{FullPath: absPath, Info: nil}, statErr
	}

	return FileInformation{FullPath: absPath, Info: stat}, nil
}
