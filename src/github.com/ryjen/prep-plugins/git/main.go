package main

import (
	"os"
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"gopkg.in/libgit2/git2go.v26"
)

func Resolve(p *plugin.Plugin) error {

	params, err := p.ReadResolver()

	if err != nil {
		return err
	}

	cloneOptions := &git.CloneOptions{
		Bare: true,
	}
	_, err = git.Clone(params.Location, params.Path, cloneOptions)

	if err != nil {
		return err
	}

	return p.WriteReturn(params.Path)
}

func main() {

	p := plugin.NewPlugin("git")

	p.OnResolve = Resolve

	err := p.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}