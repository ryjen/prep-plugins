package main

import (
	"flag"
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var TEST_DATA = flag.String("data", "", "dir of the test data")

func TestArchive(t *testing.T) {

	p := NewAutotoolsPlugin()

	params, err := support.CreateTestBuild()

	if err != nil {
		t.Error(err)
		return
	}

	params.Package = "autotools-test-plugin"
	params.Version = "0.1.0"
	params.SourcePath = filepath.Join(*TEST_DATA, "autotools")

	defer os.RemoveAll(params.RootPath)

	var Header = []string{
		"BUILD\n",
		fmt.Sprintln(params.Package),
		fmt.Sprintln(params.Version),
		fmt.Sprintln(params.SourcePath),
		fmt.Sprintln(params.BuildPath),
		fmt.Sprintln(params.InstallPath),
		fmt.Sprintln(params.BuildOpts),
		"TEST_ENV=true\n",
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

	if os.Getenv("TEST_ENV") != "true" {
	    t.Error("Did not set environment variable")
	    return
	}
}

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
