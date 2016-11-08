package config

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/mikkeloscar/flis/context"
	xkb "github.com/mikkeloscar/go-xkbcommon"
)

var configContent = []byte(`$mod: Mod1`)

func setupEnvironment(t *testing.T) (string, string) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	wd := path.Join(pwd, ".tmp")
	knownPath := path.Join(wd, ".flis")
	err = os.MkdirAll(knownPath, 0755)
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	configFile := path.Join(wd, "config.yml")
	err = ioutil.WriteFile(configFile, configContent, 0644)
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	configFileKnown := path.Join(knownPath, "config")
	err = ioutil.WriteFile(configFileKnown, configContent, 0644)
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	err = os.Setenv("HOME", wd)
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	err = os.Setenv("XDG_CONFIG_HOME", "")
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	return configFile, wd
}

func clearEnvironment(path string, t *testing.T) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}

func TestLoadConfig(t *testing.T) {
	// setup test environment
	configFile, wd := setupEnvironment(t)
	// load config given pathname
	_, err := LoadConfig(configFile)
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	// load config from known locations
	_, err = LoadConfig("")
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	err = os.Setenv("HOME", "/")
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	// fail to load config from known locations (HOME is /)
	_, err = LoadConfig("")
	if err == nil {
		t.Errorf("should fail")
	}

	// cleanup test environment
	clearEnvironment(wd, t)
}

// TestGetConfig tests getting config from context.
func TestGetConfig(t *testing.T) {
	ctx := context.Context(map[string]interface{}{
		"config": New(),
	})

	conf := Get(ctx)
	if conf == nil {
		t.Errorf("expected to get conf, got nil")
	}
}

func TestVarReplace(t *testing.T) {
	strs := map[string]string{
		"$mod+a":  "Mod1+a",
		"$mod1+a": "$mod1+a",
	}

	conf := New()
	conf.vars["$mod"] = "Mod1"

	for in, out := range strs {
		if conf.VarReplace(in) != out {
			t.Errorf("expecter '%s', got '%s'", out, conf.VarReplace(in))
		}
	}
}

func TestBindings(t *testing.T) {
	conf := New()
	conf.Modes["default"] = []*Binding{{}}

	if len(conf.Bindings()) != 1 {
		t.Errorf("expecter %d binding, got %d", 1, len(conf.Bindings()))
	}

	conf.mode = "none"

	if len(conf.Bindings()) != 0 {
		t.Errorf("expecter %d binding, got %d", 0, len(conf.Bindings()))
	}
}

func TestAddBinding(t *testing.T) {
	bindings := []struct {
		conf    *Config
		mode    string
		binding *Binding
		success bool
	}{
		{
			&Config{},
			"none",
			nil,
			false,
		},
		{
			&Config{
				Modes: map[string][]*Binding{
					"default": nil,
				},
			},
			"default",
			nil,
			true,
		},
		{
			&Config{
				Modes: map[string][]*Binding{
					"default": {
						{Keys: []xkb.KeySym{xkb.KeyA}},
					},
				},
			},
			"default",
			&Binding{Keys: []xkb.KeySym{xkb.KeyA}},
			true,
		},
		{
			&Config{
				Modes: map[string][]*Binding{
					"default": {
						{Keys: []xkb.KeySym{xkb.KeyA}},
						{Keys: []xkb.KeySym{xkb.KeyB}},
					},
				},
			},
			"default",
			&Binding{Keys: []xkb.KeySym{xkb.KeyB}},
			true,
		},
		{
			&Config{
				Modes: map[string][]*Binding{
					"default": {
						{Keys: []xkb.KeySym{xkb.KeyA}},
					},
				},
			},
			"default",
			&Binding{Keys: []xkb.KeySym{xkb.KeyB}},
			true,
		},
	}

	for _, x := range bindings {
		err := x.conf.AddBinding(x.mode, x.binding)
		if err != nil && x.success {
			t.Errorf("adding binding should not fail: %s", err)
		}

		if err == nil && !x.success {
			t.Errorf("adding binding should fail")
		}
	}
}
