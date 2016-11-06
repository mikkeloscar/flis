package commands

import (
	"fmt"
	"testing"

	"github.com/mikkeloscar/flise/context"
)

type mockCommand string

func (m mockCommand) Exec(ctx context.Context) error {
	return nil
}

func (m mockCommand) String() string {
	return string(m)
}

type mockErrorCommand string

func (m mockErrorCommand) Exec(ctx context.Context) error {
	return fmt.Errorf("exec failed")
}

func (m mockErrorCommand) String() string {
	return string(m)
}

func TestCommandString(t *testing.T) {
	cmds := map[string]*Command{
		"command1": {
			mockCommand("command1"),
			nil,
			nil,
		},
		"command1, command2": {
			mockCommand("command1"),
			nil,
			&Command{
				mockCommand("command2"),
				nil,
				nil,
			},
		},
	}

	for res, cmd := range cmds {
		if res != cmd.String() {
			t.Errorf("expected %s, got %s", res, cmd.String())
		}
	}
}

func TestCommandRun(t *testing.T) {
	cmds := []struct {
		valid   bool
		command *Command
	}{
		{
			true,
			&Command{
				mockCommand("command1"),
				nil,
				nil,
			},
		},
		{
			true,
			&Command{
				mockCommand("command2"),
				nil,
				&Command{
					mockCommand("command2"),
					nil,
					nil,
				},
			},
		},
		{
			false,
			&Command{
				mockErrorCommand("command1"),
				nil,
				nil,
			},
		},
	}

	for _, c := range cmds {
		err := c.command.Run(nil)
		if err != nil && c.valid {
			t.Errorf("command should not fail: %s", err)
		}

		if err == nil && !c.valid {
			t.Errorf("command should fail")
		}
	}
}

func TestCommandParse(t *testing.T) {
	cmds := map[string]bool{
		"":          false,
		",":         false,
		" ":         false,
		"nocommand": false,
		"exec cmd":  true,
		"exec":      false,
		"exit":      true,
	}

	for cmd, valid := range cmds {
		_, err := Parse(cmd)
		if err != nil && valid {
			t.Errorf("parsing '%s' should not fail, got: %s", cmd, err)
		}

		if err == nil && !valid {
			t.Errorf("parsing '%s' should not succeed", cmd)
		}
	}
}
