package main

import (
	"os"
	"io/ioutil"
	"testing"
	"fmt"
	"github.com/ryjen/prep-plugins/support"
)

const var TEST_DATA = "/Users/ryanjennings/.devel/prep/plugins/tests/data"

func TestCmake(t *testing.T) {

	p := NewCmakePlugin()

	params, err := plugin.CreateTestBuild()

	if err != nil {
		t.Error(err)
		return
	}

	params.Package = "cmake-plugin-test"
	params.Version = "0.1.0"
	params.SourcePath = filepath.Join(TEST_DATA, "cmake")

	defer os.RemoveAll(params.RootPath)

	var Header = []string {
		"BUILD\n",
		fmt.Sprintln(params.Package),
		fmt.Sprintln(params.Version),
		fmt.Sprintln(params.SourcePath),
		fmt.Sprintln(params.BuildPath),
		fmt.Sprintln(params.InstallPath),
		fmt.Sprintln(params.BuildOpts),
		"END\n",
	}

	err = p.ExecutePipe(Header)

	if err != nil {
		t.Error(err)
		return
	}

	fileInfo, err := ioutil.ReadDir(params.BuildPath)

	if err != nil {
		t.Error(err)
		return
	}

	if len(fileInfo) == 0 {
		t.Error("Did not generates build in build path")
		return
	}

	// TODO: archive specific checks
}

func TestMain( m *testing.M) {

	os.Exit(m.Run())
}
