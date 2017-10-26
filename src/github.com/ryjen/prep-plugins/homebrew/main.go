package main

import (
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"os"
)

func Load(p *support.Plugin) error {

	err := p.ExecuteExternal("brew", "--version")

	if err != nil {
		p.SetEnabled(false)
		p.WriteEcho(fmt.Sprint(p.Name, " not available, plugin disabled"))
	}
	return nil
}

func Install(p *support.Plugin) error {

	params, err := p.ReadInstall()

	if err != nil {
		return err
	}

	return p.ExecuteExternal("brew", "install", params.Package)
}

func Remove(p *support.Plugin) error {
	params, err := p.ReadInstall()

	if err != nil {
		return err
	}

	return p.ExecuteExternal("brew", "uninstall", params.Package)
}

func main() {

	p := support.NewPlugin("make")

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
