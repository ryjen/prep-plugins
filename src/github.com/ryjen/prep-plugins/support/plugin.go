package plugin

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
    OnLoad func(p *Plugin) error
    OnBuild func(p *Plugin) error
    OnInstall func(p *Plugin) error
    OnRemove func(p *Plugin) error
    OnResolve func(p *Plugin) error
}

/**
 * base struct that defines a package parameters
 */
type PackageParams struct {
    Package string
    Version string
}

/** 
 * build hook parameters
 */
type BuildParams struct {
    *PackageParams
    SourcePath string
    BuildPath string
    InstallPath string
    BuildOpts string
}

/** 
 * install/remove hook parameters
 */
type InstallParams struct {
    *PackageParams
    Repository string
}

/**
 * resolve hook params
 */
type ResolverParams struct {
    Path string
    Location string
}

/** 
 * creates a new plugin with no-op callbacks
 */
func New() *Plugin {
    noop := func(p *Plugin) error {
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

    params.Package, err = p.Read()

    if err != nil {
        return params, err
    }

    params.Version, err = p.Read()

    if err != nil {
        return params, err
    }

    params.SourcePath, err = p.Read()

    if err != nil {
        return params, err
    }

    params.BuildPath, err = p.Read()

    if err != nil {
        return params, err
    }

    params.InstallPath, err = p.Read()

    if err != nil {
        return params, err
    }

    params.BuildOpts, err = p.Read()

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

    params.Package, err = p.Read()

    if err != nil {
        return params, err
    }

    params.Version, err = p.Read()

    if err != nil {
        return params, err
    }

    params.Repository, err = p.Read()

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

    params.Path, err = p.Read()

    if err != nil {
        return params, err
    }

    params.Location, err = p.Read()

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
        return p.OnLoad(p)
    case "INSTALL":
        return p.OnInstall(p)
    case "REMOVE":
        return p.OnRemove(p)
    case "BUILD":
        return p.OnBuild(p)
    case "RESOLVE":
        return p.OnResolve(p)
    default:
        return errors.New("unknown plugin hook")
    }
}
