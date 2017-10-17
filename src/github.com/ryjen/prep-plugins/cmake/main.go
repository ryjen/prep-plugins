package main

import (
	"os"
	"os/exec"
	"fmt"
	"github.com/ryjen/prep-plugins/support"
)

func MakeBuild(p *plugin.Plugin) error {

	params, err := p.ReadBuild()

	if err != nil {
		return err
	}

	os.Chdir(params.BuildPath)

	cmd := exec.Command("cmake", fmt.Sprint("-DCMAKE_INSTALL_PREFIX=", params.InstallPath), params.BuildOpts, params.SourcePath)

	return cmd.Run()
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