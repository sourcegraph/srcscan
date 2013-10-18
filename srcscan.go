package srcscan

import (
	"go/build"
	"os"
	"path/filepath"
)

// Config specifies options for Scan.
type Config struct {
	// Base is the base directory that all source unit paths are made relative to. Paths within the
	// concrete source unit structs are relative to the source unit path, not Base. If Base is the
	// empty string, the current working directory is used.
	Base string

	// Profiles is the list of profiles to use when scanning for source units. If nil,
	// AllProfiles is used.
	Profiles []Profile

	// SkipDirs is a list of names of directories that are skipped while scanning.
	SkipDirs []string

	// PathIndependent, if true, indicates that all filesystem paths should be relativized, if
	// possible, or else cleared.
	PathIndependent bool

	NodeJSPackage NodeJSPackageConfig
	GoPackage     GoPackageConfig
	Ruby          RubyConfig
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
	SkipDirs: []string{"node_modules", "vendor", "testdata", "site-packages", "bower_components"},
	NodeJSPackage: NodeJSPackageConfig{
		TestDirs:          []string{"test", "tests", "spec", "specs", "unit", "mocha", "karma", "testdata"},
		TestSuffixes:      []string{"test.js", "tests.js", "spec.js", "specs.js"},
		SupportDirs:       []string{"build_support"},
		SupportFilenames:  []string{"Gruntfile.js", "build.js", "Makefile.dryice.js", "build.config.js"},
		ExampleDirs:       []string{"example", "examples", "sample", "samples", "doc", "docs", "demo", "demos"},
		ScriptDirs:        []string{"bin", "script", "scripts", "tool", "tools"},
		GeneratedDirs:     []string{"build", "dist"},
		GeneratedSuffixes: []string{".min.js", "-min.js", ".optimized.js", "-optimized.js"},
		VendorDirs:        []string{"vendor", "bower_components", "node_modules", "assets", "public", "static", "resources"},
	},
	GoPackage: GoPackageConfig{
		BuildContext: build.Default,
	},
	Ruby: RubyConfig{
		TestDirs:   []string{"spec", "specs", "test", "tests"},
		GemSrcDirs: []string{"lib"},
		AppSrcDirs: []string{"app", "lib", "config", "db"},
	},
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

	c.Base, _ = filepath.Abs(c.Base)

	for _, profile := range profiles {
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

				if profile.Dir != nil && profile.Dir.DirMatches(path, filenames) {
					relpath, abspath := c.relAbsPath(path)
					found = append(found, profile.Unit(abspath, relpath, c, info))
					if profile.TopLevelOnly {
						return filepath.SkipDir
					}
				}
			} else {
				if profile.File != nil && profile.File.FileMatches(path) {
					relpath, abspath := c.relAbsPath(path)
					found = append(found, profile.Unit(abspath, relpath, c, info))
				}
			}
			return
		})
	}

	return
}

func (c Config) relAbsPath(path string) (rel string, abs string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}

	rel, err = filepath.Rel(c.Base, abs)
	if err != nil {
		panic(err.Error())
	}
	return
}
