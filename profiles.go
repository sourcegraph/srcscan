package srcscan

import (
	"go/build"
	"strings"
)

// Profile represents criteria for a source unit and instructions on how to create it. For example,
// a simple profile might represent the following: "if there is a directory with a package.json
// file, designate the directory as a node.js package and include all *.js files in its
// subdirectories."
type Profile struct {
	// Name describes the source unit matched by this Profile.
	Name string

	Dir DirMatcher

	Unit func(path string, config Config) Unit
}

func (p Profile) DirMatches(path string, filenames []string) bool {
	if p.Dir.DirMatches(path, filenames) {
		return true
	}
	return false
}

type FileMatcher interface {
	FileMatches(path string) bool
}

type DirMatcher interface {
	DirMatches(path string, filenames []string) bool
}

// FileInDir matches directories containing a file with the specified name.
type FileInDir struct{ Filename string }

func (c FileInDir) DirMatches(path string, filenames []string) bool {
	for _, f := range filenames {
		if f == c.Filename {
			return true
		}
	}
	return false
}

// FileSuffixInDir matches directories containing a file with the specified filename suffix.
type FileSuffixInDir struct{ Suffix string }

func (c FileSuffixInDir) DirMatches(path string, filenames []string) bool {
	for _, f := range filenames {
		if strings.HasSuffix(f, c.Suffix) {
			return true
		}
	}
	return false
}

var AllProfiles = []Profile{
	Profile{
		Name: "node.js package",
		Dir:  FileInDir{"package.json"},
		Unit: func(dir string, config Config) Unit {
			u := &NodeJSPackage{Dir: dir}
			u.read(config)
			return u
		},
	},
	Profile{
		Name: "Python package",
		Dir:  FileInDir{"__init__.py"},
		Unit: func(dir string, config Config) Unit {
			return &PythonPackage{dir}
		},
	},
	Profile{
		Name: "Go package",
		Dir:  FileSuffixInDir{".go"},
		Unit: func(dir string, config Config) Unit {
			u := &GoPackage{Package: build.Package{Dir: dir}}
			u.read(config)
			return u
		},
	},
}
