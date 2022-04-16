package givematlib

import (
	"os"
	"path/filepath"
)

func InDataDir(folders ...string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	base := []string{homeDir, ".local", "share", "givematerial"}
	p := append(base, folders...)
	path := filepath.Join(p...)
	return path, nil
}
