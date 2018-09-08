package config

import (
	"fmt"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"
	wlc "github.com/mikkeloscar/go-wlc"
	xkb "github.com/mikkeloscar/go-xkbcommon"

	"gopkg.in/yaml.v2"
)

// Map xkb modifier names to the equivalent wlc modifier bitmasks.
var modifiers = map[string]wlc.ModifierBit{
	xkb.ModNameShift: wlc.BitModShift,
	xkb.ModNameCaps:  wlc.BitModCaps,
	xkb.ModNameCtrl:  wlc.BitModCtrl,
	"Ctrl":           wlc.BitModCtrl,
	xkb.ModNameAlt:   wlc.BitModAlt,
	"Alt":            wlc.BitModAlt,
	xkb.ModNameNum:   wlc.BitModMod2,
	"Mod3":           wlc.BitModMod3,
	xkb.ModNameLogo:  wlc.BitModLogo,
	"Mod5":           wlc.BitModMod5,
}

// yaml config structure.
type yamlConfig struct {
	Gaps struct {
		Inner uint
		Outer uint
	}
	Bindings []string
	Modes    []struct {
		Name     string
		Bindings []string
	}
	Outputs []struct {
		ID         string
		Enabled    bool
		Size       string
		Pos        string
		Background string
	}
}

// parseConfig parses config from content.
func parseConfig(content []byte) (*Config, error) {
	t := yamlConfig{}
	c := make(map[interface{}]interface{})

	err := yaml.Unmarshal(content, &t)
	if err != nil {
		return nil, err
	}

	// this will not error since we already managed to unmarshal the
	// content above.
	_ = yaml.Unmarshal(content, c)

	config := New()

	// extract variables
	for k, v := range c {
		switch key := k.(type) {
		case string:
			if isVar(key) {
				config.vars[key] = v.(string)
				continue
			}
		}
	}

	// parse default bindings
	err = parseBindings(config, "default", t.Bindings)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bindings: %s", err)
	}

	// parse modes
	for _, mode := range t.Modes {
		if _, ok := config.Modes[mode.Name]; !ok {
			config.Modes[mode.Name] = make([]*Binding, 0, len(mode.Bindings))
		}

		err = parseBindings(config, mode.Name, mode.Bindings)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bindings: %s", err)
		}
	}

	return config, nil
}

func parseBindings(config *Config, mode string, bindings []string) error {
	for _, binding := range bindings {
		b, err := parseBinding(config, binding)
		if err != nil {
			log.Warn(err)
			continue
		}

		err = config.AddBinding(mode, b)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseBinding(config *Config, bindingStr string) (*Binding, error) {
	var err error
	replaced := config.VarReplace(bindingStr)

	split := strings.SplitN(replaced, " ", 2)

	if len(split) != 2 {
		return nil, fmt.Errorf("invalid binding: '%s'", bindingStr)
	}

	keys := strings.Split(split[0], "+")

	binding := &Binding{Raw: bindingStr}

	for _, key := range keys {
		if mod, ok := modifiers[key]; ok {
			binding.Modifiers |= uint32(mod)
			continue
		}

		sym := xkb.KeySymFromName(key, xkb.KeySymCaseInsensitive)
		if sym == xkb.KeyNoSymbol {
			return nil, fmt.Errorf("unknown key '%s' in binding '%s'", key, replaced)
		}

		binding.Keys = append(binding.Keys, sym)
	}

	binding.Command, err = cmdParse(split[1], config)
	if err != nil {
		return nil, err
	}

	return binding, nil
}

// check if value is a variable.
func isVar(value string) bool {
	if !(len(value) > 1 && value[0] == '$') {
		return false
	}

	for _, c := range value[1:] {
		if !isAlphaNumericUnderscore(c) {
			return false
		}
	}

	return true
}

// isAlphaNumericUnderscore returns true if rune is lower case character,
// number or underscore.
func isAlphaNumericUnderscore(r rune) bool {
	return r == '_' || unicode.IsDigit(r) || (unicode.IsLetter(r) && unicode.IsLower(r))
}
