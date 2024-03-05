package copier

import (
	"io/fs"
	"time"
)

type destFileInfo struct {
	FileName    string
	FileSize    int64
	IsDirectory bool
}

func MakeFakeDestinationFileInfo(fileName string, fileSize int64) destFileInfo {
	return destFileInfo{FileName: fileName, FileSize: fileSize}
}

func (dest destFileInfo) Name() string {
	return dest.FileName
}

func (dest destFileInfo) Size() int64 {
	return dest.FileSize
}

func (destFileInfo) Mode() fs.FileMode {
	return 0
}

func (destFileInfo) ModTime() time.Time {
	return time.Now()
}

func (dest destFileInfo) IsDir() bool {
	return dest.IsDirectory
}

func (destFileInfo) Sys() any {
	return nil
}
