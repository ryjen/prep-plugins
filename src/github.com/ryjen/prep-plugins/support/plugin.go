package plugin

import (
    "bufio"
    "os"
    "strings"
    "fmt"
    "errors"
    "os/exec"
    "io"
    "io/ioutil"
    "net/http"
    "path/filepath"
    "mime"
    "github.com/mholt/archiver"
)

type Hook func(*Plugin) error

/**
 * the plugin type with hook callbacks
 */
type Plugin struct {
    Name string
    OnLoad Hook
    OnUnload Hook
    OnBuild Hook
    OnInstall Hook
    OnRemove Hook
    OnResolve Hook
    Input io.Reader
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
    PackageParams
    SourcePath string
    BuildPath string
    InstallPath string
    BuildOpts string
}

/** 
 * install/remove hook parameters
 */
type InstallParams struct {
    PackageParams
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
func NewPlugin(name string) *Plugin {
    noop := func(p *Plugin) error {
        return nil
    }

    return &Plugin{name, noop, noop,noop,
    noop, noop, noop, os.Stdin }
}

/**
 * creates new build hook parameters
 */
func NewBuildParams() *BuildParams {
    return new(BuildParams)
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

func (p *Plugin) keyName(key string) string {

    keys := []string{"PREP", strings.ToUpper(p.Name), strings.ToUpper(key)}

    return strings.Join(keys, "_")
}

/**
 * abstraction for plugin key/value storage
 * just uses environment variables for now
 */
func (p *Plugin) Save(key string, value string) error {

    return os.Setenv(p.keyName(key), value)
}

func (p *Plugin) Lookup(key string) string {
    return os.Getenv(p.keyName(key))
}


func (p *Plugin) SetEnabled(value bool) error {
    if value {
        return p.Save("enabled", "yes")
    } else {
        return p.Save("enabled", "no")
    }
}

func (p *Plugin) IsEnabled() bool {
    return p.Lookup("enabled") != "no"
}

/**
 * read a line and trim it
 */
func (p *Plugin) Read() (string, error) {

    reader := bufio.NewReader(p.Input)
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
    if !p.IsEnabled() {
        return nil
    }

    command, err := p.Read()

    if err != nil {
        if err == io.EOF {
            return nil
        }
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
    case "UNLOAD":
        return p.OnUnload(p)
    default:
        return errors.New(fmt.Sprint("unknown plugin hook ", command))
    }
}

func (plugin *Plugin) ExecutePipe(header []string) error {

    reader, writer := io.Pipe()

    defer reader.Close()

    // set the plugin input to our pipe
    plugin.Input = reader

    handler := make(chan error)

    go func() {
        defer writer.Close()

        for _, data := range header {

            if _, err := writer.Write([]byte(data)); err != nil {
                handler <- err
                return
            }
        }

        close(handler)
    }()

    err := plugin.Execute()

    if err != nil {
        return err
    }

    return <- handler
}

func (p *Plugin) ExecuteExternal(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

/**
 * build parameters for testing
 */
type TestBuildParams struct {
    BuildParams
    RootPath string
}


type FileVersionInfo struct {
    FileName string
    Version string
    BaseName string
}

/**
 * utility to parse a version from a file path
 * @return a version string (ex. 4.3.2-beta)
 */
func parseFileAndVerionFromPath(path string) (*FileVersionInfo, error) {

    info := &FileVersionInfo{}

    // get the filename
    info.FileName = filepath.Base(path)

    // split by extension separator
    parts := strings.Split(info.FileName, "-")

    var version []string

    if len(parts) > 1 {
        info.BaseName = parts[0]
        parts = strings.Split(strings.Join(parts[1:], "-"), ".")
    } else {
        parts = strings.Split(info.FileName, ".")
    }

    // while a valid mime type extension
    for len(parts) > 0 {

        // pop the last extension part again    
        part := parts[0]

        // retest if valid mime type extension
        mtype := mime.TypeByExtension("." + part)

        if (len(mtype) > 0) {
            break;
        }
        
        parts = parts[1:]

        version = append(version, part)
    }

    // and join the rest to get the version string
    info.Version = strings.Join(version, ".")

    return info, nil
}

/**
 * creates parameters for testing build plugins.
 * It is up to the test to remove the temporary RootPath
 */
func CreateTestBuild(url string) (*TestBuildParams, error) {

    params := &TestBuildParams{}

    var err error = nil

    params.RootPath, err = ioutil.TempDir(os.TempDir(), filepath.Base(os.TempDir()))

    if err != nil {
        return nil, err
    }

    info, err := parseFileAndVerionFromPath(url)

    if err != nil {
        return nil, err
    }

    params.Package, params.Version = info.BaseName, info.Version

    sourceFolder := params.Package + "-" + params.Version

    archiveFile := info.FileName

    resp, err := http.Get(url)

    if err != nil {
        return nil, err
    }

    archivePath := filepath.Join(params.RootPath, archiveFile)

    file, err := os.Create(archivePath)

    if err != nil {
        return nil, err
    }

    if resp.StatusCode != 200 {
        return nil, errors.New("Invalid url")
    }

    _, err = io.Copy(file, resp.Body)

    if err != nil {
        return nil, err
    }

    file.Close()

    params.SourcePath = filepath.Join(params.RootPath, "source")

    err = os.MkdirAll(params.SourcePath, os.FileMode(0700))

    if err != nil {
        return nil, err
    }

    params.BuildPath = filepath.Join(params.RootPath, "build")

    err = os.MkdirAll(params.BuildPath, os.FileMode(0700))

    if err != nil {
        return nil, err
    }

    params.InstallPath = filepath.Join(params.RootPath, "install")

    err = os.MkdirAll(params.InstallPath, os.FileMode(0700))

    if err != nil {
        return nil, err
    }

    ar := archiver.MatchingFormat(archiveFile)

    err = ar.Open(archivePath, params.SourcePath)

    if err != nil {
        return nil, err
    }

    params.SourcePath = filepath.Join(params.SourcePath, sourceFolder)

    return params, nil
}

