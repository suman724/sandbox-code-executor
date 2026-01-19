package workspace

import (
	"errors"
	"os"
)

func SecureDelete(path string) error {
	if path == "" {
		return errors.New("missing path")
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		return os.RemoveAll(path)
	}
	if err := overwriteFile(path, info.Size()); err != nil {
		return err
	}
	return os.Remove(path)
}

func overwriteFile(path string, size int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	zero := make([]byte, 4096)
	remaining := size
	for remaining > 0 {
		chunk := int64(len(zero))
		if remaining < chunk {
			chunk = remaining
		}
		if _, err := file.Write(zero[:chunk]); err != nil {
			return err
		}
		remaining -= chunk
	}
	return file.Sync()
}
