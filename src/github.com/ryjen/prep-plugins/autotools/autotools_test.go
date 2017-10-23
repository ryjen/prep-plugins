package main

import (
	"os"
	"io/ioutil"
	"testing"
	"fmt"
)

func TestArchive(t *testing.T) {

	p := NewAutotoolsPlugin()

	path, err := ioutil.TempDir(os.TempDir(), t.Name())

	if err != nil {
		t.Error("unable to create temporary folder ", path)
	}

	filename := "libarchive-3.3.2.tar.gz"

	url := strings.Join("http://www.libarchive.org/downloads/", filename)

	resp, err := http.Get(url)

	if err != nil {
		t.Error(err)
	}

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return err
	}

	file.Close()

	buildpath := filepath.Join(path, "build")

	err = os.MkdirAll(buildpath)

	if err != nil {
		return err
	}

	installpath := filepath.Join(path, "install")

	err = os.MkdirAll(installpath)

	if err != nil {
		return err
	}

	ar := archiver.MatchingFormat(filename)

	err = ar.Open(path, params.Path)

	if err != nil {
		return err
	}

	var Header = []string{
		"BUILD\n",
		"libarchive\n",
		"3.3.2\n",
		fmt.Sprintln(path),
		fmt.Sprintln(buildpath),
		fmt.Sprintln(installpath),
		"\n",
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
}

func TestMain( m *testing.M) {

	os.Exit(m.Run())
}