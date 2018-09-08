package config

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/mikkeloscar/flis/backend"
)

// Exit implements the exit command.
type Exit struct{}

// Exec executes the exit command.
func (e Exit) Exec(ctx context.Context) error {
	// TODO: close all views (maybe in compositor terminate cb)
	log.Infof("Terminating compositor...")
	backend.Get(ctx).Terminate()
	return nil
}

// String returns a string formatting of the exit command.
func (e Exit) String() string {
	return "exit"
}

// parseExec parses an exit command definition.
func parseExit(lex *lexer, config *Config) (Executer, error) {
	return Exit{}, nil
}
