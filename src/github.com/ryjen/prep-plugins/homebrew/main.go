package main

import (
	"os"
	"github.com/ryjen/prep-plugins/support"
	"fmt"
)

func Load(p *plugin.Plugin) error {

	err := p.RunCommand("brew", "--version")

	if err != nil {
		p.SetEnabled(false)
		p.WriteEcho(fmt.Sprint(p.Name, " not available, plugin disabled"))
	}
	return nil
}

func Install(p *plugin.Plugin) error {

	params, err := p.ReadInstall()

	if err != nil {
		return err
	}

	return p.RunCommand("brew", "install", params.Package)
}

func Remove(p *plugin.Plugin) error {
	params, err := p.ReadInstall()

	if err != nil {
		return err
	}

	return p.RunCommand("brew", "uninstall", params.Package)
}

func main() {

	p := plugin.NewPlugin("make")

	p.OnLoad = Load
	p.OnInstall = Install
	p.OnRemove = Remove

	err := p.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}