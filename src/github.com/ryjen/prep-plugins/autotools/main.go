package main

import (
	"os"
	"fmt"
	"errors"
	"github.com/ryjen/prep-plugins/support"
)

func Load(p *plugin.Plugin) error {

	err := p.ExecuteExternal("autoconf", "--version")

	if err != nil {
		p.SetEnabled(false)
		p.WriteEcho(fmt.Sprint(p.Name, " not available, plugin disabled"))
	}
	return nil
}


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
		err = p.ExecuteExternal(autogen)

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
	return p.ExecuteExternal(configure, fmt.Sprint("--prefix=", params.InstallPath), params.BuildOpts)
}

func NewAutotoolsPlugin() *plugin.Plugin {

	p := plugin.NewPlugin("autotools")

	p.OnLoad = Load
	p.OnBuild = MakeBuild

	return p
}

func main() {

	err := NewAutotoolsPlugin().Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}