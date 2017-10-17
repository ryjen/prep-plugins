package main

import (
	"os"
	"github.com/ryjen/prep-plugins/support"
	"gopkg.in/libgit2/git2go.v26"
)

func MakeBuild(p *plugin.Plugin) error {

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

	p := plugin.New()

	p.OnBuild = MakeBuild

	err := p.Execute()

	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}