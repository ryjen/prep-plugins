package main

import (
	"github.com/ryjen/prep-plugins/support"
	"os"
)

func Load(p *support.Plugin) error {

	err := p.ExecuteQuiet("apt-get", "--version")

	if err != nil {
		return support.NotFoundError(err)
	}
	return nil
}

func Add(p *support.Plugin) error {

	params, err := p.ReadAddRemove()

	if err != nil {
		return err
	}

	err = os.Chdir(os.TempDir())

	if err != nil {
		return err
	}

	err = p.ExecuteExternal("apt-get", "download", params.Package)

	if err != nil {
		return err
	}

	return p.ExecuteExternal("dpkg", "--force-not-root", "--root=\""+params.Repository+"\"", "-i", params.Package+"*.deb")
}

func Remove(p *support.Plugin) error {
	params, err := p.ReadAddRemove()

	if err != nil {
		return err
	}

	err = os.Chdir(os.TempDir())

	if err != nil {
		return err
	}

	return p.ExecuteExternal("dpkg", "-i", params.Package+"*.deb", "--force-not-root", "--root=\""+params.Repository+"\"")
}

func main() {

	p := support.NewPlugin("apt")

	p.OnLoad = Load
	p.OnAdd = Add
	p.OnRemove = Remove

	err := p.Execute()

	if err != nil {
		os.Exit(support.ErrorCode(err))
	}

	os.Exit(0)
}
