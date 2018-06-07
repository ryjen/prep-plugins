package support

import (
	"C"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type Hook func(*Plugin) error

/**
 * the plugin type with hook callbacks
 */
type Plugin struct {
	Name      string
	OnLoad    Hook
	OnUnload  Hook
	OnBuild   Hook
	OnTest    Hook
	OnInstall Hook
	OnAdd     Hook
	OnRemove  Hook
	OnResolve Hook
	Input     io.Reader
	Output    io.Writer
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
	SourcePath  string
	BuildPath   string
	InstallPath string
	BuildOpts   string
}

type BuiltParams struct {
	PackageParams
	SourcePath string
	BuildPath  string
}

/**
 * add/remove hook parameters
 */
type AddRemoveParams struct {
	PackageParams
	Repository string
}

/**
 * resolve hook params
 */
type ResolverParams struct {
	Path     string
	Location string
}

/**
 * creates a new plugin with no-op callbacks
 */
func NewPlugin(name string) *Plugin {
	noop := func(p *Plugin) error {
		return nil
	}

	return &Plugin{name, noop, noop, noop,
		noop, noop, noop, noop, noop, os.Stdin, os.Stdout}
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

	if err != nil {
		return err
	}

	for !strings.EqualFold(line, "END") {

		env := strings.Split(line, "=")

		if len(env) == 2 {
			os.Setenv(strings.TrimSpace(env[0]), strings.TrimSpace(env[1]))
		}

		line, err = p.Read()

		if err != nil {
			return err
		}
	}

	return err
}

/**
 * read build hook parameters
 */
func (p *Plugin) ReadBuild() (*BuildParams, error) {

	params := new(BuildParams)

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
 * read build hook parameters
 */
func (p *Plugin) ReadBuilt() (*BuiltParams, error) {

	params := new(BuiltParams)

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

	err = p.ReadEnvVars()

	return params, err
}

/**
 * read install hook parameters
 */
func (p *Plugin) ReadAddRemove() (*AddRemoveParams, error) {
	params := new(AddRemoveParams)

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
	params := new(ResolverParams)

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
	_, err := fmt.Fprintln(p.Output, "RETURN", value)
	return err
}

/**
 * write an echo message
 */
func (p *Plugin) WriteEcho(value string) error {
	_, err := fmt.Fprintln(p.Output, "ECHO", value)
	return err
}

/**
 * reads an input hook, and executes
 */
func (p *Plugin) Execute() error {
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
	case "ADD":
		return p.OnAdd(p)
	case "REMOVE":
		return p.OnRemove(p)
	case "BUILD":
		return p.OnBuild(p)
	case "TEST":
		return p.OnTest(p)
	case "INSTALL":
		return p.OnInstall(p)
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

	return <-handler
}

func (p *Plugin) ExecuteExternal(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *Plugin) ExecuteQuiet(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (p *Plugin) ExecuteOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	b, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func GetErrorCode(err error) int {
	// try to get the exit code
	if exitError, ok := err.(*exec.ExitError); ok {
		ws := exitError.Sys().(syscall.WaitStatus)
		return ws.ExitStatus()
	}
	return -1
}

/**
 * build parameters for testing
 */
type TestBuildParams struct {
	BuildParams
	InstallPath string
	RootPath    string
}

/**
 * creates parameters for testing build plugins.
 * It is up to the test to remove the temporary RootPath
 */
func CreateTestBuild() (*TestBuildParams, error) {

	params := &TestBuildParams{}

	var err error = nil

	params.RootPath, err = ioutil.TempDir(os.TempDir(), filepath.Base(os.TempDir()))

	if err != nil {
		return nil, err
	}

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

	return params, nil
}

// export copy
func Copy(src, dst string) (int64, error) {

	if len(src) == 0 || len(dst) == 0 {
		return 0, nil
	}

	// stat the source file
	stat, err := os.Stat(src)

	var srcReader io.ReadCloser

	// source is not a file?
	if err != nil && os.IsNotExist(err) {

		// try a url
		resp, err := http.Get(src)

		if err != nil {
			return 0, fmt.Errorf("%s is not a regular file or url", src)
		}

		srcReader = resp.Body

	} else if stat.Mode().IsRegular() {

		// check to copy the file name if not specified
		stat, err = os.Stat(dst)
		if stat != nil && stat.Mode().IsDir() {
			dst = filepath.Join(dst, filepath.Base(src))
		}

		// open the source file
		srcReader, err = os.Open(src)
		if err != nil {
			return 0, err
		}
	} else {
		// possibly a dir
		return 0, err
	}

	defer srcReader.Close()

	destFile, err := os.Create(dst)

	if err != nil {
		return 0, err
	}

	defer destFile.Close()

	return io.Copy(destFile, srcReader)
}
