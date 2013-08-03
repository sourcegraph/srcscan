package srcscan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Unit represents a "source unit," such as a Go package, a node.js package, or a Python package.
type Unit interface {
	Path() string
}

type DirUnit struct {
	Dir string
}

// Path implements Unit.
func (d DirUnit) Path() string {
	return d.Dir
}

// Units implements sort.Interface.
type Units []Unit

func (u Units) Len() int      { return len(u) }
func (u Units) Swap(i, j int) { u[i], u[j] = u[j], u[i] }
func (u Units) Less(i, j int) bool {
	return fmt.Sprintf("%T", u[i])+u[i].Path() < fmt.Sprintf("%T", u[j])+u[j].Path()
}

// NodeJSPackage represents a node.js package.
type NodeJSPackage struct {
	DirUnit
	PackageJSON    json.RawMessage `json:",omitempty"`
	LibFiles       []string        `json:",omitempty"`
	ScriptFiles    []string        `json:",omitempty"`
	SupportFiles   []string        `json:",omitempty"`
	ExampleFiles   []string        `json:",omitempty"`
	TestFiles      []string        `json:",omitempty"`
	VendorFiles    []string        `json:",omitempty"`
	GeneratedFiles []string        `json:",omitempty"`
}

type NodeJSPackageConfig struct {
	TestDirs          []string
	TestSuffixes      []string
	SupportDirs       []string
	SupportFilenames  []string
	ExampleDirs       []string
	ScriptDirs        []string
	GeneratedDirs     []string
	GeneratedSuffixes []string
	VendorDirs        []string
}

func (u *NodeJSPackage) read(config Config) {
	// Read package.json.
	var err error
	u.PackageJSON, err = ioutil.ReadFile(filepath.Join(u.Dir, "package.json"))
	if err != nil {
		panic("read package.json: " + err.Error())
	}

	// Populate *Files fields.
	c := config.NodeJSPackage
	err = filepath.Walk(u.Dir, func(path string, info os.FileInfo, inerr error) (err error) {
		if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".js") {
			relpath, _ := filepath.Rel(u.Dir, path)
			parts := strings.Split(relpath, "/")
			for _, part := range parts {
				if contains(c.VendorDirs, part) {
					u.VendorFiles = append(u.VendorFiles, relpath)
					return
				} else if contains(c.GeneratedDirs, part) || hasAnySuffix(c.GeneratedSuffixes, relpath) {
					u.GeneratedFiles = append(u.GeneratedFiles, relpath)
					return
				} else if contains(c.ScriptDirs, part) {
					u.ScriptFiles = append(u.ScriptFiles, relpath)
					return
				} else if contains(c.ExampleDirs, part) {
					u.ExampleFiles = append(u.ExampleFiles, relpath)
					return
				} else if contains(c.TestDirs, part) || hasAnySuffix(c.TestSuffixes, relpath) {
					u.TestFiles = append(u.TestFiles, relpath)
					return
				} else if contains(c.SupportDirs, part) || contains(c.SupportFilenames, info.Name()) {
					u.SupportFiles = append(u.SupportFiles, relpath)
					return
				}
			}
			u.LibFiles = append(u.LibFiles, relpath)
		} else if info.IsDir() {
			if info.Name() == "node_modules" {
				return filepath.SkipDir
			}

			// Don't traverse into sub-packages.
			if path != u.Dir && dirHasFile(path, "package.json") {
				return filepath.SkipDir
			}
		}
		return
	})
	if err != nil {
		panic("scan files: " + err.Error())
	}
}

// GoPackage represents a Go package.
type GoPackage struct {
	DirUnit
}

// PythonPackage represents a Python package.
type PythonPackage struct {
	DirUnit
}
