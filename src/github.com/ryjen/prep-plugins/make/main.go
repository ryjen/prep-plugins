package main

import (
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"os"
	"path/filepath"
)

func Load(p *support.Plugin) error {

	err := p.ExecuteExternal("make", "--version")

	if err != nil {
		p.SetEnabled(false)
		p.WriteEcho(fmt.Sprint(p.Name, " not available, plugin disabled"))
	}
	return nil
}

func MakeBuild(p *support.Plugin) error {

	params, err := p.ReadBuild()

	if err != nil {
		return err
	}

	os.Chdir(params.BuildPath)

	return p.ExecuteExternal("make", "-f", filepath.Join(params.SourcePath, "Makefile"), "-I", params.SourcePath, "install")
}

func NewMakePlugin() *support.Plugin {

	p := support.NewPlugin("make")

	p.OnLoad = Load
	p.OnBuild = MakeBuild

	return p
}

func main() {

	err := NewMakePlugin().Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
