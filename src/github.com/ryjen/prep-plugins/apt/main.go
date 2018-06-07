package main

import (
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"os"
)

func Load(p *support.Plugin) error {

	err := p.ExecuteQuiet("apt-get", "--version")

	if err != nil {
		p.SetEnabled(false)
		p.WriteEcho(fmt.Sprint(p.Name, " not available, plugin disabled"))
	}
	return nil
}

func Add(p *support.Plugin) error {

	params, err := p.ReadAddRemove()

	if err != nil {
		return err
	}

	err = p.ExecuteExternal("apt-cache", "show", params.Package)

	if err == nil {
		err = p.ExecuteExternal("apt-get", "install", params.Package)
	}

	return err
}

func Remove(p *support.Plugin) error {
	params, err := p.ReadAddRemove()

	if err != nil {
		return err
	}

	return p.ExecuteExternal("apt-get", "remove", params.Package)
}

func main() {

	p := support.NewPlugin("apt")

	p.OnLoad = Load
	p.OnAdd = Add
	p.OnRemove = Remove

	err := p.Execute()

	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
