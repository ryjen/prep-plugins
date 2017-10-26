package main

import (
    "flag"
	"os"
	"io/ioutil"
	"testing"
	"path/filepath"
	"fmt"
)

var TEST_DATA = flag.String("data", "", "dir of the test data")

func TestArchive(t *testing.T) {

	p := NewArchivePlugin()

	path, err := ioutil.TempDir(os.TempDir(), t.Name())

	if err != nil {
		t.Error("unable to create temporary folder ", path)
		return
	}

	var archiveFile = filepath.Join(*TEST_DATA, "archive", "autotools.tar.gz")

	var Header = []string{
		"RESOLVE\n",
		fmt.Sprintln(path),
		fmt.Sprintln(archiveFile),
		"END\n",
	}

	err = p.ExecutePipe(Header)

	if err != nil {
		t.Error(err)
		return
	}

	fileInfo, err := ioutil.ReadDir(path)

	if err != nil {
		t.Error(err)
		return
	}

	if len(fileInfo) == 0 {
		t.Error("Did not extract test archive")
	}

	// TODO: archive specific checks
}

func TestMain( m *testing.M) {
	os.Exit(m.Run())
}