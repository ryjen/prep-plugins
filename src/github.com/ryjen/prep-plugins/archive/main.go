package main

import (
	"github.com/mholt/archiver"
	"github.com/ryjen/prep-plugins/support"
	"os"
	"path/filepath"
	"strings"
	"mime"
	"errors"
)

type FileVersionInfo struct {
	FileName string
	Version  string
	BaseName string
}

/**
 * utility to parse a version from a file path
 * @return a version string (ex. 4.3.2-beta)
 */
func parseFileAndVerionFromPath(path string) (*FileVersionInfo, error) {

	info := &FileVersionInfo{}

	// get the filename
	info.FileName = filepath.Base(path)

	// split by extension separator
	parts := strings.Split(info.FileName, "-")

	var version []string

	if len(parts) > 1 {
		info.BaseName = parts[0]
		parts = strings.Split(strings.Join(parts[1:], "-"), ".")
	} else {
		parts = strings.Split(info.FileName, ".")
	}

	// while a valid mime type extension
	for len(parts) > 0 {

		// pop the last extension part again
		part := parts[0]

		// retest if valid mime type extension
		mtype := mime.TypeByExtension("." + part)

		if len(mtype) > 0 {
			break
		}

		parts = parts[1:]

		version = append(version, part)
	}

	// and join the rest to get the version string
	info.Version = strings.Join(version, ".")

	return info, nil
}

func Resolve(p *support.Plugin) error {

	params, err := p.ReadResolver()

    if len(params.Location) == 0 || len(params.Path) == 0 {
        return errors.New("invalid parameter")
    }

	if err != nil {
		return err
	}

	err = os.MkdirAll(params.Path, os.FileMode(755))

	if err != nil {
		return err
	}

	filename := filepath.Base(params.Location)

	path := filepath.Join(os.TempDir(), filename)

	_, err = support.Copy(params.Location, path)

	if err != nil {
		return err
	}

	ar := archiver.MatchingFormat(filename)

	err = ar.Open(path, params.Path)

	if err != nil {
		return err
	}

	info, err := parseFileAndVerionFromPath(filename)

	if info != nil {
		newPath := filepath.Join(params.Path, info.BaseName)

		stat, _ := os.Stat(newPath)

		if stat != nil && stat.Mode().IsDir() {
			params.Path = newPath
		}
	}

	return p.WriteReturn(params.Path)
}

func NewArchivePlugin() *support.Plugin {

	p := support.NewPlugin("archive")

	p.OnResolve = Resolve

	return p
}

func main() {

	err := NewArchivePlugin().Execute()

	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
