package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/mikkeloscar/go-wlc"
)

// Config defines the configuration of a compositor.
type Config struct {
	vars       map[string]string
	mode       string
	Modes      map[string][]*Binding
	Workspaces map[uint]string
	Bars       map[string]Bar
	Outputs    map[string]Output
	Gaps       struct {
		Inner uint
		Outer uint
	}
	// Inputs map[string]Input
}

// Get config from context.
func Get(ctx context.Context) *Config {
	return ctx.Value("config").(*Config)
}

// New initializes a new default config.
func New() *Config {
	return &Config{
		vars: make(map[string]string),
		mode: "default",
		Modes: map[string][]*Binding{
			"default": make([]*Binding, 0),
		},
		Workspaces: make(map[uint]string, 0),
		Bars:       make(map[string]Bar),
		Outputs:    make(map[string]Output),
		Gaps: struct {
			Inner uint
			Outer uint
		}{0, 0},
	}
}

// LoadConfig loads config file from disk.
func LoadConfig(path string) (*Config, error) {
	content, err := loadConfig(path)
	if err != nil {
		return nil, err
	}

	return parseConfig(content)
}

// loadConfig tries to load config file from the path specified by the function
// fargument or by a list of default locations.
func loadConfig(file string) ([]byte, error) {
	if file != "" {
		return ioutil.ReadFile(file)
	}

	// create list of paths to look for config file.
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = path.Join(os.Getenv("HOME"), ".config")
	}

	paths := []string{
		path.Join(os.Getenv("HOME"), ".flis/config"),
		path.Join(xdgConfigHome, "flis/config"),
		// TODO sysconfdir,
	}

	for _, file := range paths {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Warnf("unable to load config: %s", err)
			continue
		}
		return content, nil
	}

	return nil, fmt.Errorf("failed to load config file")
}

var varRegexp = regexp.MustCompile(`\$\w+`)

// VarReplace replaces variables with the actual value.
// TODO: implement this  better
func (c *Config) VarReplace(str string) string {
	vars := varRegexp.FindAllString(str, -1)
	for _, v := range vars {
		if val, ok := c.vars[v]; ok {
			str = strings.Replace(str, v, val, -1)
		}
	}

	return str
}

// Bindings returns a list of bindings for the current mode.
func (c *Config) Bindings() []*Binding {
	if b, ok := c.Modes[c.mode]; ok {
		return b
	}

	return nil
}

// AddBinding adds a binding to the current mode.
func (c *Config) AddBinding(mode string, b *Binding) error {
	bindings, ok := c.Modes[mode]
	if !ok {
		return fmt.Errorf("failed to find %s mode in config", mode)
	}

	if len(bindings) == 0 {
		c.Modes[mode] = []*Binding{b}
		return nil
	}

	found := false
	for i, binding := range bindings {
		// if equal overwrite previous binding.
		// TODO: debug log this?
		if b.Equal(binding) {
			bindings[i] = b
			// bindings = append(bindings[:i-1], bindings[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		bindings = append(bindings, b)
	}

	sort.Sort(Bindings(bindings))
	c.Modes[mode] = bindings
	return nil
}

// Bar defines the configuration of a bar.
type Bar struct {
	ID string
}

// Output defines the configuration of an output.
type Output struct {
	ID         string
	Enabled    bool
	Size       wlc.Size
	Pos        wlc.Point
	Background struct {
		Path string
		Mode string
	}
}
