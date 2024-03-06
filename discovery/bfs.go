package discovery

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type bfsDiscoverer struct {
	State BfsState
}

type BfsState struct {
	Queue []FileInformation
}

func MakeBfsDiscoverer(basePath string) (Discoverer, error) {
	baseInfo, baseInfoErr := MakeFileInformation(basePath)
	if baseInfoErr != nil {
		return nil, baseInfoErr
	}

	initialState := BfsState{Queue: []FileInformation{baseInfo}}

	if !baseInfo.Info.IsDir() && !baseInfo.Info.Mode().IsRegular() {
		return nil, fmt.Errorf("%v is an unsupported file", basePath)
	}

	return &bfsDiscoverer{State: initialState}, nil
}

func (bfs *bfsDiscoverer) Next() (FileInformation, error) {
	if len(bfs.State.Queue) == 0 {
		return FileInformation{}, io.EOF
	}

	var resultItem *FileInformation
	for resultItem == nil && len(bfs.State.Queue) > 0 {
		currentItem := bfs.State.Queue[0]
		canProcessItem := currentItem.Info.IsDir() || currentItem.Info.Mode().IsRegular()
		bfs.State.Queue = bfs.State.Queue[1:]
		if canProcessItem {
			resultItem = &currentItem
		}
	}

	if resultItem == nil {
		return FileInformation{}, io.EOF
	}

	if resultItem.Info.IsDir() {
		childern, _ := os.ReadDir(resultItem.FullPath)

		for _, child := range childern {
			fInfo, fInfoErr := MakeFileInformation(filepath.Join(resultItem.FullPath, child.Name()))
			if fInfoErr == nil {
				bfs.State.Queue = append(bfs.State.Queue, fInfo)
			}
		}
	}

	return *resultItem, nil
}
