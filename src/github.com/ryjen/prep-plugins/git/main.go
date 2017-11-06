package main

import (
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"os"
	"errors"
	"strings"
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

	if len(params.Location) == 0 || len(params.Path) == 0 {
	    return errors.New("invalid parameter")
	}

	stat, err := os.Stat(params.Path)

	if stat != nil && stat.IsDir() {

        err = os.Chdir(params.Path)

        if err != nil {
            return err
        }

        origin, err := p.ExecuteOutput("git", "remote", "get-url", "origin")

        if err != nil {
            return err
        }

        if strings.TrimSpace(origin) == params.Location {
            return p.WriteReturn(params.Path)
        }
    }

	fmt.Println("Cloning ", params.Location, " to ", params.Path)

	err = p.ExecuteExternal("git", "clone", params.Location, params.Path)

	if err != nil {
		return err
	}

	err = os.Chdir(params.Path)

	if err != nil {
		return err
	}

	err = p.ExecuteExternal("git", "submodule", "update", "--init", "--recursive")

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
