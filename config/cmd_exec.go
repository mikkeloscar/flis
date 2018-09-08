package config

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/mikkeloscar/flis/backend"
)

// Exec implements the exec command.
type Exec string

// Exec executes the exec command.
func (e Exec) Exec(ctx context.Context) error {
	log.Debugf("Executing shell command '%s'", e)
	backend.Get(ctx).Exec("/bin/sh", "-c", string(e))
	return nil
}

// String returns a string formatting of the exec command.
func (e Exec) String() string {
	return string(e)
}

// parseExec parses an exec command definition.
func parseExec(lex *lexer, config *Config) (Executer, error) {
	var args []string

	for t := lex.nextItem(); t.typ == itemString; t = lex.nextItem() {
		if strings.ContainsRune(t.val, ' ') {
			t.val = fmt.Sprintf(`"%s"`, t.val)
		}
		args = append(args, t.val)
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no command defined")
	}

	return Exec(strings.Join(args, " ")), nil
}
