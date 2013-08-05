package srcscan

import (
	"github.com/kr/pretty"
	"go/build"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	Default.PathIndependent = true
	type scanTest struct {
		config *Config
		dir    string
		units  []Unit
	}
	tests := []scanTest{
		{
			dir: "testdata",
			units: []Unit{
				&GoPackage{
					Package: build.Package{
						Dir:            "testdata/go",
						Name:           "mypkg",
						ImportPath:     "github.com/sourcegraph/srcscan/testdata/go",
						GoFiles:        []string{"a.go", "b.go"},
						Imports:        []string{},
						ImportPos:      nil,
						TestGoFiles:    []string{"a_test.go"},
						TestImports:    []string{},
						TestImportPos:  nil,
						XTestGoFiles:   []string{"b_test.go"},
						XTestImports:   []string{},
						XTestImportPos: nil,
					},
				},
				&GoPackage{
					Package: build.Package{
						Dir:            "testdata/go/cmd/mycmd",
						Name:           "main",
						ImportPath:     "github.com/sourcegraph/srcscan/testdata/go/cmd/mycmd",
						GoFiles:        []string{"mycmd.go"},
						Imports:        []string{},
						ImportPos:      nil,
						TestGoFiles:    nil,
						TestImports:    []string{},
						TestImportPos:  nil,
						XTestGoFiles:   nil,
						XTestImports:   []string{},
						XTestImportPos: nil,
					},
				},
				&GoPackage{
					Package: build.Package{
						Dir:            "testdata/go/qux",
						Name:           "qux",
						ImportPath:     "github.com/sourcegraph/srcscan/testdata/go/qux",
						GoFiles:        []string{"qux.go"},
						Imports:        []string{},
						ImportPos:      nil,
						TestGoFiles:    nil,
						TestImports:    []string{},
						TestImportPos:  nil,
						XTestGoFiles:   nil,
						XTestImports:   []string{},
						XTestImportPos: nil,
					},
				},
				&NodeJSPackage{
					DirUnit:        DirUnit{Dir: "testdata/node.js"},
					PackageJSON:    []byte(`{"name":"mypkg"}`),
					LibFiles:       []string{"a.js", "lib/a.js"},
					TestFiles:      []string{"a_test.js", "test/b.js", "test/c_test.js"},
					VendorFiles:    []string{"vendor/a.js"},
					GeneratedFiles: []string{"a.min.js", "dist/a.js"},
				},
				&NodeJSPackage{
					DirUnit:     DirUnit{Dir: "testdata/node.js/subpkg"},
					PackageJSON: []byte(`{"name":"subpkg"}`),
					LibFiles:    []string{"a.js"},
				},
				&PythonPackage{DirUnit{"testdata/python/mypkg"}},
				&PythonPackage{DirUnit{"testdata/python/mypkg/qux"}},
			},
		},
	}
	for _, test := range tests {
		// Use default config if config is nil.
		var config Config
		if test.config != nil {
			config = *test.config
		} else {
			config = Default
		}

		units, err := config.Scan(test.dir)
		if err != nil {
			t.Errorf("got error %q", err)
			continue
		}

		sort.Sort(Units(units))
		sort.Sort(Units(test.units))
		if !reflect.DeepEqual(test.units, units) {
			t.Errorf("units:\n%v", pretty.Diff(test.units, units))
			if len(test.units) == len(units) {
				for i := range test.units {
					if !reflect.DeepEqual(test.units[i], units[i]) {
						t.Errorf("units[%d]:\n%v", i, strings.Join(pretty.Diff(test.units[i], units[i]), "\n"))
					}
				}
			}
		}
	}
}
