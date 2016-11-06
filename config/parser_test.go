package config

import "testing"

func TestParseConfig(t *testing.T) {
	tests := []struct {
		content []byte
		success bool
	}{
		{
			[]byte(`$mod: Mod1`),
			true,
		},
		// invalid yaml
		{
			[]byte(`$mod`),
			false,
		},
		// invalid binding (no error just warn log)
		{
			[]byte(`
bindings:
  - Mod1+a`),
			true,
		},
		{
			[]byte(`
modes:
  - name: resize
    bindings:
      - Mod1+q exit`),
			true,
		},
	}

	for _, test := range tests {
		_, err := parseConfig(test.content)
		if err != nil && test.success {
			t.Errorf("parsing config should not fail: %s", err)
		}

		if err == nil && !test.success {
			t.Errorf("parsing config should fail")
		}
	}
}

func TestParseBindings(t *testing.T) {
	tests := []struct {
		config   *Config
		mode     string
		bindings []string
		success  bool
	}{
		{
			&Config{
				Modes: map[string][]*Binding{
					"default": []*Binding{},
				},
				mode: "default",
				vars: map[string]string{
					"$mod": "Mod1",
				},
			},
			"default",
			[]string{"$mod+A exec command", "Mod1+a exit"},
			true,
		},
		{
			&Config{
				Modes: map[string][]*Binding{
					"default": []*Binding{},
				},
				mode: "default",
				vars: map[string]string{
					"$mod": "Mod1",
				},
			},
			"default",
			[]string{"$mod+A", "Mod1+InvalidKey exit", "Mod1+A exec"},
			true,
		},
		{
			&Config{
				Modes: map[string][]*Binding{
					"default": []*Binding{},
				},
				mode: "default",
				vars: map[string]string{
					"$mod": "Mod1",
				},
			},
			"none",
			[]string{"$mod+A exec command"},
			false,
		},
	}

	for _, test := range tests {
		err := parseBindings(test.config, test.mode, test.bindings)
		if err != nil && test.success {
			t.Errorf("parse bindings should not fail: %s", err)
		}

		if err == nil && !test.success {
			t.Errorf("parse bindings should fail")
		}
	}
}

func TestIsVar(t *testing.T) {
	vars := map[string]bool{
		"$mod":   true,
		"$a1_":   true,
		"$ab_12": true,
		"$AB":    false,
		"$<":     false,
		"abc":    false,
	}

	for v, ret := range vars {
		if ret != isVar(v) {
			t.Errorf("expected '%s' to be %t, got %t", v, ret, isVar(v))
		}
	}
}
