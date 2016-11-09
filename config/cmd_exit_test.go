package config

import (
	"testing"

	"github.com/mikkeloscar/flis/backend"
	"github.com/mikkeloscar/flis/context"
)

func TestExitRun(t *testing.T) {
	ctx := context.Context(map[string]interface{}{
		"backend": backend.Mock{},
	})
	exit := Exit{}

	err := exit.Exec(ctx)
	if err != nil {
		t.Errorf("failed to execute exit command: %s", err)
	}
}

func TestExitString(t *testing.T) {
	cmd := "exit"
	exit := Exit{}

	if exit.String() != cmd {
		t.Errorf("expected %s, got %s", cmd, exit.String())
	}
}

func TestParseExit(t *testing.T) {
	_, err := parseExit(lex(""), nil)
	if err != nil {
		t.Errorf("parseExit should not fail, got: %s", err)
	}
}
