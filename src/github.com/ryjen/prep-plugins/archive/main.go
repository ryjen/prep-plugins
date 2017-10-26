package main

import (
	"os"
	"fmt"
	"net/http"
	"path/filepath"
	"github.com/ryjen/prep-plugins/support"
	"github.com/mholt/archiver"
	"io"
	"strings"
	"net/url"
)

func Resolve(p *plugin.Plugin) error {

	params, err := p.ReadResolver()

	if err != nil {
		return err
	}

	err = os.MkdirAll(params.Path, os.FileMode(755))

	if err != nil {
		return err
	}

	filename := filepath.Base(params.Location)

	path := filepath.Join(os.TempDir(), filename)

	url, err := url.Parse(params.Location)

	if strings.Contains(url.Scheme, "http") {

		resp, err := http.Get(params.Location)

		if err != nil {
			return err
		}

		file, err := os.Create(path)

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(file, resp.Body)

		if err != nil {
			return err
		}
	} else {
		plugin.Copy(params.Location, path)
	}

	ar := archiver.MatchingFormat(filename)

	err = ar.Open(path, params.Path)

	if err != nil {
		return err
	}

	return p.WriteEcho(params.Path)
}

func NewArchivePlugin() *plugin.Plugin {

	p := plugin.NewPlugin("archive")

	p.OnResolve = Resolve

	return p
}

func main() {

	err := NewArchivePlugin().Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}