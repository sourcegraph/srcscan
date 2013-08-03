package srcscan

import (
	"os"
	"path/filepath"
	"strings"
)

func contains(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func hasAnySuffix(suffixes []string, str string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix(str, s) {
			return true
		}
	}
	return false
}

func dirHasFile(dir, filename string) bool {
	path := filepath.Join(dir, filename)
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	} else if err == nil && info.Mode().IsRegular() {
		return true
	}
	panic("dirHasFile: " + err.Error())
}
