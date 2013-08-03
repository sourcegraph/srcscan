package srcscan

import (
	"os"
	"path/filepath"
)

// Config specifies options for Scan.
type Config struct {
	// Profiles is the list of profiles to use when scanning for source units. If nil,
	// AllProfiles is used.
	Profiles []Profile

	// SkipDirs is a list of names of directories that are skipped while scanning.
	SkipDirs []string
}

func (c Config) skipDir(name string) bool {
	for _, dirToSkip := range c.SkipDirs {
		if name == dirToSkip {
			return true
		}
	}
	return false
}

var Default = Config{
	SkipDirs: []string{"node_modules", "vendor"},
}

// Scan is shorthand for Default.Scan.
func Scan(dir string) (found []Unit, err error) {
	return Default.Scan(dir)
}

// Scan walks the directory tree at dir, looking for source units that match profiles in the
// configuration. Scan returns a list of all source units found.
func (c Config) Scan(dir string) (found []Unit, err error) {
	var profiles []Profile
	if c.Profiles != nil {
		profiles = c.Profiles
	} else {
		profiles = AllProfiles
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, inerr error) (err error) {
		if inerr != nil {
			return inerr
		}
		if info.IsDir() {
			if dir != path && c.skipDir(info.Name()) {
				return filepath.SkipDir
			}

			var dirh *os.File
			dirh, err = os.Open(path)
			if err != nil {
				return
			}
			defer dirh.Close()

			var filenames []string
			filenames, err = dirh.Readdirnames(0)
			if err != nil {
				return
			}
			for _, p := range profiles {
				if p.Dir != nil && p.Dir.DirMatches(path, filenames) {
					found = append(found, p.Unit(path))
				}
			}
		}
		return
	})

	return
}
