package main

import (
	"github.com/ryjen/prep-plugins/support"
	"os"
	"path/filepath"
)

func Load(p *support.Plugin) error {

	err := p.ExecuteQuiet("make", "--version")

	if err != nil {
		return support.NotFoundError(err)
	}
	return nil
}

func MakeBuildArg(p *support.Plugin, sourcePath string, buildPath string, arg string) error {

	err := os.Chdir(buildPath)

	if err != nil {
		return err
	}

	// try to make from build path first
	file := filepath.Join(buildPath, "Makefile")

	_, err = os.Stat(file)

	if len(arg) == 0 {
		if err == nil {
			err = p.ExecuteExternal("make", "-f", file, "-I", buildPath)
		} else {
			err = p.ExecuteExternal("make", "-f", filepath.Join(sourcePath, "Makefile"), "-I", sourcePath)
		}
	} else {
		if err == nil {
			err = p.ExecuteExternal("make", "-f", file, "-I", buildPath, arg)
		} else {
			err = p.ExecuteExternal("make", "-f", filepath.Join(sourcePath, "Makefile"), "-I", sourcePath, arg)
		}
	}

	return err
}

func MakeBuild(p *support.Plugin) error {

	params, err := p.ReadBuild()

	if err != nil {
		return err
	}

	return MakeBuildArg(p, params.SourcePath, params.BuildPath, "")
}

func MakeTest(p *support.Plugin) error {

	params, err := p.ReadBuilt()

	if err != nil {
		return err
	}

	return MakeBuildArg(p, params.SourcePath, params.BuildPath, "test")
}

func MakeInstall(p *support.Plugin) error {

	params, err := p.ReadBuilt()

	if err != nil {
		return err
	}

	return MakeBuildArg(p, params.SourcePath, params.BuildPath, "install")
}

func NewMakePlugin() *support.Plugin {

	p := support.NewPlugin("make")

	p.OnLoad = Load
	p.OnBuild = MakeBuild
	p.OnTest = MakeTest
	p.OnInstall = MakeInstall

	return p
}

func main() {

	err := NewMakePlugin().Execute()

	if err != nil {
		os.Exit(support.ErrorCode(err))
	}

	os.Exit(0)
}
