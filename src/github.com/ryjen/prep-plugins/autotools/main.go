package main

import (
	"os"
	"os/exec"
	"fmt"
	"errors"
	"github.com/ryjen/prep-plugins/support"
)

func MakeBuild(p *plugin.Plugin) error {

	params, err := p.ReadBuild()

	if err != nil {
		return err
	}

	// path to the configure script
	configure := fmt.Sprint(params.BuildPath, "/configure")

	_, err = os.Stat(configure)

	// if configure script doesn't exist...
	if err != nil && os.IsNotExist(err) {

		// path to the autogen script
		autogen := fmt.Sprint(params.SourcePath, "/autogen.sh")

		_, err = os.Stat(autogen)

		// if doesn't exist we don't know how to build
		if err != nil && os.IsNotExist(err) {
			return errors.New(fmt.Sprint("Don't know how to build ", params.Package, "... no autotools configuration found."))
		}

		// change to the source path
		err = os.Chdir(params.SourcePath)

		if err != nil {
			return err
		}

		// run the autogen script
		cmd := exec.Command(autogen)

		err = cmd.Run()

		if err != nil {
			return err
		}

		// test the configure script again
		_, err = os.Stat(configure)

		// it should exist now
		if err != nil && os.IsNotExist(err) {
			return errors.New(fmt.Sprint("Could not generate a configure script for ", params.Package))
		}
	}

	// go to the build path
	err = os.Chdir(params.BuildPath)

	if err != nil {
		return err
	}

	// and execute the configure script
	cmd := exec.Command(configure, fmt.Sprint("--prefix=", params.InstallPath), params.BuildOpts)

	return cmd.Run()
}

func main() {

	p := plugin.New()

	p.OnBuild = MakeBuild

	err := p.Execute()

	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}