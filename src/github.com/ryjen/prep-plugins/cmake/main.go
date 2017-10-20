package main

import (
	"os"
	"fmt"
	"github.com/ryjen/prep-plugins/support"
)

func Load(p *plugin.Plugin) error {

	err := p.ExecuteExternal("cmake", "--version")

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

	os.Chdir(params.BuildPath)

	return p.ExecuteExternal("cmake", fmt.Sprint("-DCMAKE_INSTALL_PREFIX=", params.InstallPath), params.BuildOpts, params.SourcePath)
}

func main() {

	p := plugin.NewPlugin("cmake")

	p.OnLoad = Load
	p.OnBuild = MakeBuild

	err := p.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}