package srcscan

import (
	"github.com/kr/pretty"
	"reflect"
	"sort"
	"testing"
)

func TestScan(t *testing.T) {
	type scanTest struct {
		config *Config
		dir    string
		units  []Unit
	}
	tests := []scanTest{
		{
			dir: "testdata",
			units: []Unit{
				&GoPackage{DirUnit{"testdata/go/cmd/mycmd"}},
				&GoPackage{DirUnit{"testdata/go/qux"}},
				&GoPackage{DirUnit{"testdata/go"}},
				&NodeJSPackage{DirUnit{"testdata/node.js"}},
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
		}
	}
}
