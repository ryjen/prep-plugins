package main

import (
	"flag"
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"os"
	"path/filepath"
	"testing"
)

var TEST_DATA = flag.String("data", "", "dir of the test data")

func TestMake(t *testing.T) {

	p := NewMakePlugin()

	params, err := support.CreateTestBuild()

	if err != nil {
		t.Error(err)
		return
	}

	params.Package = "make-plugin-test"
	params.Version = "0.1.0"

	params.SourcePath = filepath.Join(*TEST_DATA, "make")

	defer os.RemoveAll(params.RootPath)

	_, err = support.Copy(filepath.Join(params.SourcePath, "Makefile"), params.BuildPath)

	if err != nil {
		t.Error(err)
		return
	}

	var Header = []string{
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

	path := filepath.Join(params.BuildPath, "make_output.log")

	_, err = os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		t.Error(err)
		return
	}
}

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
