package logic

import (
	"os"
)

func IsDirectory(name string) (bool, error) {
	statRes, err := os.Stat(name)
	if err != nil {
		return false, err
	}

	return statRes.IsDir(), nil
}
