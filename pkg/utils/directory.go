package utils

import (
	"errors"
	"os"
)

func EnsureDirectory(path string) error {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return os.Mkdir(path, 0755)
		}

		return err
	}

	return nil
}
