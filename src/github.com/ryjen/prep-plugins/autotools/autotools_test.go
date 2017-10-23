package main

import (
	"os"
	"io/ioutil"
	"testing"
	"fmt"
	"io"
	"path/filepath"
	"net/http"
	"github.com/mholt/archiver"
	"github.com/ryjen/prep-plugins/support"
)

func CreateBuildDirectories(t *testing.T) (string, *plugin.BuildParams) {

	params := plugin.NewBuildParams()

	path, err := ioutil.TempDir(os.TempDir(), t.Name())

	if err != nil {
		t.Error(err)
	}

	sourceFolder := "libarchive-3.3.2"

	archiveFile := sourceFolder + ".tar.gz"

	url := "http://www.libarchive.org/downloads/" + archiveFile

	resp, err := http.Get(url)

	if err != nil {
		t.Error(err)
	}

	archivePath := filepath.Join(path, archiveFile)

	file, err := os.Create(archivePath)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Error("Invalid url")
	}

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		t.Error(err)
	}

	file.Close()

	params.Package = "libarchive"
	params.Version = "3.3.2"

	params.SourcePath = filepath.Join(path, "source")

	err = os.MkdirAll(params.SourcePath, os.FileMode(0700))

	if err != nil {
		t.Error(err)
	}

	params.BuildPath = filepath.Join(path, "build")

	err = os.MkdirAll(params.BuildPath, os.FileMode(0700))

	if err != nil {
		t.Error(err)
	}

	params.InstallPath = filepath.Join(path, "install")

	err = os.MkdirAll(params.InstallPath, os.FileMode(0700))

	if err != nil {
		t.Error(err)
	}

	ar := archiver.MatchingFormat(archiveFile)

	err = ar.Open(archivePath, params.SourcePath)

	if err != nil {
		t.Error(err)
	}

	params.SourcePath = filepath.Join(params.SourcePath, sourceFolder)

	return path, params
}


func TestArchive(t *testing.T) {

	p := NewAutotoolsPlugin()

	path, params := CreateBuildDirectories(t)

	defer os.RemoveAll(path)

	fmt.Println(params.SourcePath)

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

	err := p.ExecutePipe(Header)

	if err != nil {
		t.Error(err)
	}

	fileInfo, err := ioutil.ReadDir(params.BuildPath)

	if err != nil {
		t.Error(err)
	}

	if len(fileInfo) == 0 {
		t.Error("Did not generates build in build path")
	}
}

func TestMain( m *testing.M) {

	os.Exit(m.Run())
}