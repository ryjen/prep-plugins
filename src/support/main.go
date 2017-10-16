package main

import (
    "bufio"
    "os"
    "strings"
    "fmt"
    "errors"
)

/**
 * the plugin type with hook callbacks
 */
type Plugin struct {
    OnLoad func() error
    OnBuild func() error
    OnInstall func() error
    OnRemove func() error
    OnResolve func() error
}

/**
 * base struct that defines a package parameters
 */
type PackageParams struct {
    pkg string
    version string
}

/** 
 * build hook parameters
 */
type BuildParams struct {
    *PackageParams
    sourcePath string
    buildPath string
    installPath string
    buildOpts string
}

/** 
 * install/remove hook parameters
 */
type InstallParams struct {
    *PackageParams
    repository string
}

/**
 * resolve hook params
 */
type ResolverParams struct {
    path string
    location string
}

/** 
 * creates a new plugin with no-op callbacks
 */
func NewPlugin() *Plugin {
    noop := func() error {
        return nil
    }

    return &Plugin{noop, noop, noop, noop, noop }
}

/**
 * creates new build hook parameters
 */
func NewBuildParams() *BuildParams {
    return &BuildParams{}
}

/**
 * creates new install hook parameters
 */
func NewInstallParams() *InstallParams {
    return &InstallParams{}
}

/**
 * creates new resolve hook parameters
 */
func NewResolverParams() *ResolverParams {
    return &ResolverParams{}
}

/**
 * read a line and trim it
 */
func (p *Plugin) Read() (string, error) {

    reader := bufio.NewReader(os.Stdin)
    text, err := reader.ReadString('\n')

    return strings.TrimSpace(text), err
}

/**
 * read environment variables until END of header
 */
func (p *Plugin) ReadEnvVars() error {
    line, err := p.Read()

    for !strings.EqualFold(line, "END") {
        if err != nil {
            return err
        }

        env := strings.Split(line, "=")

        if len(env) == 2 {
            os.Setenv(strings.TrimSpace(env[0]), strings.TrimSpace(env[1]))
        }

        line, err = p.Read()
    }

    return err
}

/**
 * read build hook parameters
 */
func (p *Plugin) ReadBuild() (*BuildParams, error) {

    params := NewBuildParams()

    var err error

    params.pkg, err = p.Read()

    if err != nil {
        return params, err
    }

    params.version, err = p.Read()

    if err != nil {
        return params, err
    }

    params.sourcePath, err = p.Read()

    if err != nil {
        return params, err
    }

    params.buildPath, err = p.Read()

    if err != nil {
        return params, err
    }

    params.installPath, err = p.Read()

    if err != nil {
        return params, err
    }

    params.buildOpts, err = p.Read()

    if err != nil {
        return params, err
    }

    err = p.ReadEnvVars()

    return params, err
}

/**
 * read install hook parameters
 */
func (p *Plugin) ReadInstall() (*InstallParams, error) {
    params := NewInstallParams()

    var err error

    params.pkg, err = p.Read()

    if err != nil {
        return params, err
    }

    params.version, err = p.Read()

    if err != nil {
        return params, err
    }

    params.repository, err = p.Read()

    if err != nil {
        return params, err
    }

    err = p.ReadEnvVars()

    return params, err
}

/**
 * read resolver hook parameters
 */
func (p *Plugin) ReadResolver() (*ResolverParams, error) {
    params := NewResolverParams()

    var err error

    params.path, err = p.Read()

    if err != nil {
        return params, err
    }

    params.location, err = p.Read()

    if err != nil {
        return params, err
    }

    err = p.ReadEnvVars()

    return params, err
}

/**
 * write a return value
 */
func (p *Plugin) WriteReturn(value string) error {
    _, err := fmt.Println("RETURN ", value)
    return err
}

/**
 * write an echo message
 */
func (p *Plugin) WriteEcho(value string) error {
    _, err := fmt.Println("ECHO ", value)
    return err
}

/**
 * reads an input hook, and executes
 */
func  (p *Plugin) Execute() error {
    command, err := p.Read()

    if err != nil {
        return err
    }

    switch strings.ToUpper(command) {
    case "LOAD":
        return p.OnLoad()
    case "INSTALL":
        return p.OnInstall()
    case "REMOVE":
        return p.OnRemove()
    case "BUILD":
        return p.OnBuild()
    case "RESOLVE":
        return p.OnResolve()
    default:
        return errors.New("unknown plugin hook")
    }
}


func main() {

    plugin := NewPlugin()

    err := plugin.Execute()

    if err != nil {
        os.Exit(1)
    }

    os.Exit(0)
}