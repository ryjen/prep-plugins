package main

import (
	"errors"
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"os"
	"path/filepath"
)

func Load(p *support.Plugin) error {

	err := p.ExecuteQuiet("autoconf", "--version")

	if err != nil {
	    return support.NotFoundError(err)
	}

	err = p.ExecuteQuiet("automake", "--version")

	if err != nil {
		return support.NotFoundError(err)
	}

	return nil
}

func RunAutogen(p *support.Plugin, params *support.BuildParams) error {
	// path to the autogen script
	autogen := filepath.Join(params.SourcePath, "autogen.sh")

	_, err := os.Stat(autogen)

	// if doesn't exist we don't know how to build
	if err != nil && os.IsNotExist(err) {
		return errors.New(fmt.Sprint("Don't know how to build ", params.Package, "... no autotools configuration found."))
	}

	// go to the build path
	err = os.Chdir(params.SourcePath)

	if err != nil {
		return err
	}

	// run the autogen script
	return p.ExecuteExternal(autogen, params.BuildOpts, params.SourcePath)
}

func MakeBuild(p *support.Plugin) error {

	params, err := p.ReadBuild()

	if err != nil {
		return err
	}

	// path to the configure script
	configure := filepath.Join(params.SourcePath, "configure")

	_, err = os.Stat(configure)

	// if configure script doesn't exist...
	if err != nil && os.IsNotExist(err) {

		err = RunAutogen(p, params)

		if err != nil {
			return err
		}

		_, err = os.Stat(configure)

		if err != nil && os.IsNotExist(err) {
			return errors.New("Unable to generate configure script for build")
		}
	}

	// go to the build path
	err = os.Chdir(params.BuildPath)

	if err != nil {
		return err
	}

	// and execute the configure script
	return p.ExecuteExternal(configure, fmt.Sprint("--prefix=", params.InstallPath), params.BuildOpts, params.SourcePath)
}

func NewAutotoolsPlugin() *support.Plugin {

	p := support.NewPlugin("autotools")

	p.OnLoad = Load
	p.OnBuild = MakeBuild

	return p
}

func main() {

	err := NewAutotoolsPlugin().Execute()

	if err != nil {
		os.Exit(support.ErrorCode(err))
	}

	os.Exit(0)
}
