package srcscan

import (
	"os"
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

	File FileMatcher

	TopLevelOnly bool

	Unit func(abspath, relpath string, config Config, info os.FileInfo) Unit
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

type FileHasSuffix struct{ Suffix string }

func (c FileHasSuffix) FileMatches(path string) bool {
	return strings.HasSuffix(path, c.Suffix)
}

var AllProfiles = []Profile{
	Profile{
		Name: "NPM package",
		Dir:  FileInDir{"package.json"},
		Unit: readNPMPackage,
	},
	Profile{
		Name: "Bower component",
		Dir:  FileInDir{"bower.json"},
		Unit: readBowerComponent,
	},
	Profile{
		Name:         "Python package and module",
		TopLevelOnly: true,
		Dir:          FileInDir{"__init__.py"},
		File:         FileHasSuffix{".py"},
		Unit: func(abspath, relpath string, config Config, info os.FileInfo) Unit {
			if info.IsDir() {
				return &PythonPackage{relpath}
			} else {
				return &PythonModule{relpath}
			}
		},
	},
	Profile{
		Name: "Go package",
		Dir:  FileSuffixInDir{".go"},
		Unit: readGoPackage,
	},
	Profile{
		Name: "Java Maven project",
		Dir:  FileInDir{"pom.xml"},
		Unit: readJavaMavenProject,
	},
	Profile{
		Name: "Ruby Gem",
		Dir:  FileSuffixInDir{".gemspec"},
		Unit: readRubyGem,
	},
	Profile{
		Name: "Ruby app",
		Dir:  FileInDir{"config.ru"},
		Unit: readRubyApp,
	},
	// TODO(sqs): support Ruby apps (i.e., non-gem Ruby projects)
}
