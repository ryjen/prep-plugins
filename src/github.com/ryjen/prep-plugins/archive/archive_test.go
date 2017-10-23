package main

import (
	"os"
	"io/ioutil"
	"testing"
	"fmt"
)

func TestArchive(t *testing.T) {

	p := NewArchivePlugin()

	path, err := ioutil.TempDir(os.TempDir(), t.Name())

	if err != nil {
		t.Error("unable to create temporary folder ", path)
	}

	var Header = []string{
		"RESOLVE\n",
		fmt.Sprintln(path),
		"http://www.libarchive.org/downloads/libarchive-3.3.2.tar.gz\n",
		"END\n",
	}

	err = p.ExecutePipe(Header)

	if err != nil {
		t.Error(err)
	}

	fileInfo, err := ioutil.ReadDir(path)

	if err != nil {
		t.Error(err)
	}

	if len(fileInfo) == 0 {
		t.Error("Did not extract test archive")
	}

	// TODO: archive specific checks
}

func TestMain( m *testing.M) {

	os.Exit(m.Run())
}