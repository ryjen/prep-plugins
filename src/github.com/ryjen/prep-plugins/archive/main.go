package main

import (
	"os"
	"fmt"
	"github.com/ryjen/prep-plugins/support"
	"github.com/VictorLowther/go-libarchive"
)

func Resolve(p *plugin.Plugin) error {

	_, err := p.ReadResolver()

	if err != nil {
		return err
	}

	return nil
}

func main() {

	p := plugin.NewPlugin("archive")

	p.OnResolve = Resolve

	err := p.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}