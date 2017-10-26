package main

import (
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"os"
)

func Load(p *support.Plugin) error {

	err := p.ExecuteExternal("git", "--version")

	if err != nil {
		p.SetEnabled(false)
		p.WriteEcho(fmt.Sprint(p.Name, " not available, plugin disabled"))
	}
	return nil
}

func Resolve(p *support.Plugin) error {

	params, err := p.ReadResolver()

	if err != nil {
		return err
	}

	err = p.ExecuteExternal("git", "clone", params.Location, params.Path)

	if err != nil {
		return err
	}

	return p.WriteReturn(params.Path)
}

func NewGitPlugin() *support.Plugin {

	p := support.NewPlugin("git")

	p.OnLoad = Load
	p.OnResolve = Resolve

	return p
}

func main() {

	err := NewGitPlugin().Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
