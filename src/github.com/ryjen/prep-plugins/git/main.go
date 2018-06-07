package main

import (
	"errors"
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"os"
	"strings"
)

func Load(p *support.Plugin) error {

	err := p.ExecuteQuiet("git", "--version")

	if err != nil {
		p.SetEnabled(false)
		p.WriteEcho(fmt.Sprint(p.Name, " not available, plugin disabled"))
	}
	return nil
}

func IsGitError(err error) bool {
	return err != nil && support.GetErrorCode(err) == 128
}

func ResolveExisting(p *support.Plugin, params *support.ResolverParams, branch string) error {

	err := os.Chdir(params.Path)

	if err != nil {
		return err
	}

	origin, err := p.ExecuteOutput("git", "remote", "get-url", "origin")

	if err != nil {
		fmt.Println("git remote get-url error")
		return err
	}

	if strings.TrimSpace(origin) != params.Location {
		return errors.New(fmt.Sprintln(err, "Unknown origin", origin, "for", params.Location))
	}

	curr, err := p.ExecuteOutput("git", "rev-parse", "--abbrev-ref", "HEAD")

	if err != nil {
		return err
	}

	if strings.TrimSpace(curr) != branch {
		err = p.ExecuteExternal("git", "checkout", branch)

		if err != nil {
			return err
		}
	}

	err = p.ExecuteExternal("git", "pull", "-q", "origin", branch)

	if err != nil && !IsGitError(err) {
		return err
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

	branch := "master"

	extra := strings.Split(params.Location, "#")

	if len(extra) > 1 {
		branch = extra[1]
		params.Location = extra[0]
	}

	p.WriteEcho("Cloning " + params.Location + "#" + branch)

	err = p.ExecuteExternal("git", "clone", "-q", params.Location, "-b", branch, "--single-branch", params.Path)

	if err != nil {

		if IsGitError(err) {
			err = ResolveExisting(p, params, branch)

			if err != nil {
				return err
			}

		} else {
			return err
		}
	} else {

		err = os.Chdir(params.Path)

		if err != nil {
			return err
		}
	}

	p.WriteEcho("Updating submodules")

	err = p.ExecuteExternal("git", "submodule", "-q", "update", "--init", "--recursive")

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
		os.Exit(1)
	}

	os.Exit(0)
}
