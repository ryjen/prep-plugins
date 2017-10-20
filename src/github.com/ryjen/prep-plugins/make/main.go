package main

import (
	"os"
	"github.com/ryjen/prep-plugins/support"
	"fmt"
)

func Load(p *plugin.Plugin) error {

	err := p.ExecuteExternal("make", "--version")

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

    return p.ExecuteExternal("make", "-j2", "install")
}

func main() {
	
	p := plugin.NewPlugin("make")

	p.OnLoad = Load
	p.OnBuild = MakeBuild

	err := p.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}