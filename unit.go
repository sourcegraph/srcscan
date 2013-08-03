package srcscan

import (
	"fmt"
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
}

// GoPackage represents a Go package.
type GoPackage struct {
	DirUnit
}

// PythonPackage represents a Python package.
type PythonPackage struct {
	DirUnit
}
