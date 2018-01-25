package main

import (
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"os"
)

func Load(p *support.Plugin) error {

	err := p.ExecuteExternal("cmake", "--version")

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

	return p.ExecuteExternal("cmake", fmt.Sprint("-DCMAKE_INSTALL_PREFIX=", params.InstallPath), params.BuildOpts,
	    params.SourcePath)
}

func NewCmakePlugin() *support.Plugin {

	p := support.NewPlugin("cmake")

	p.OnLoad = Load
	p.OnBuild = MakeBuild

	return p
}

func main() {

	err := NewCmakePlugin().Execute()

	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
