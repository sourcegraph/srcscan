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

func hasSubdir(root, dir string) (rel string, ok bool) {
	if p, err := filepath.EvalSymlinks(root); err == nil {
		root = p
	}
	if p, err := filepath.EvalSymlinks(dir); err == nil {
		dir = p
	}
	const sep = string(filepath.Separator)
	root = filepath.Clean(root)
	if !strings.HasSuffix(root, sep) {
		root += sep
	}
	dir = filepath.Clean(dir)
	if !strings.HasPrefix(dir, root) {
		return "", false
	}
	return filepath.ToSlash(dir[len(root):]), true
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}
