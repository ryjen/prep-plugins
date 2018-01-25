package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"bytes"
)

var TEST_DATA = flag.String("data", "", "dir of the test data")

func TestArchive(t *testing.T) {

	p := NewArchivePlugin()

	path, err := ioutil.TempDir(os.TempDir(), t.Name())

	if err != nil {
		t.Error("unable to create temporary folder ", path)
		return
	}

	defer os.RemoveAll(path)

	var archiveFile = filepath.Join(*TEST_DATA, "archive", "test-0.1.0.tar.gz")

	var Header = []string{
		"RESOLVE\n",
		fmt.Sprintln(path),
		fmt.Sprintln(archiveFile),
		"END\n",
	}

	buffer := new(bytes.Buffer)

	p.Output = buffer

	err = p.ExecutePipe(Header)

	if err != nil {
		t.Error(err)
		return
	}

	n, err := fmt.Sscanf(buffer.String(), "RETURN %s", &path)

	if err != nil || n != 1{
		t.Error("no return value from plugin")
		return
	}

	_, err = os.Stat(filepath.Join(path, "README"))

	if err != nil {
		t.Error(err)
		return
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}