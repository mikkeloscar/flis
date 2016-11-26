package config

import (
	"context"
	"fmt"
)

// Executer is an interface describing how to execute commands.
type Executer interface {
	Exec(context.Context) error
	String() string
}

// Criteria - TODO: not implemented.
type Criteria int

// Command defines a command with criteria and possible chained sub commands.
type Command struct {
	Executer
	// TODO: criteria
	Criteria *Criteria
	Next     *Command
}

// String returns the string representation of a command.
func (c *Command) String() string {
	// TODO: criteria
	if c.Next != nil {
		return c.Executer.String() + ", " + c.Next.String()
	}
	return c.Executer.String()
}

// Run runs a command and any subsequent commands if it's chained.
func (c *Command) Run(ctx context.Context) error {
	// TODO: criteria
	err := c.Exec(ctx)
	if err != nil {
		return err
	}

	if c.Next != nil {
		return c.Next.Run(ctx)
	}

	return nil
}

type commandParser func(lex *lexer, config *Config) (Executer, error)

// table mapping command names to parse functions able to parse the command
// definitions.
var cmdParseTable = map[string]commandParser{
	"exec": parseExec,
	"exit": parseExit,
}

// cmdParse parses a command string into a command structure.
func cmdParse(commandStr string, config *Config) (*Command, error) {
	// TODO: criteria

	lexer := lex(commandStr)
	defer lexer.drain()

	cmdToken := lexer.nextItem()
	if cmdToken.typ != itemString {
		return nil, fmt.Errorf("expected string, got token '%s'", cmdToken.val)
	}

	// TODO: chained commands
	command := &Command{}

	fn, ok := cmdParseTable[cmdToken.val]

	if !ok {
		return nil, fmt.Errorf("command '%s' not implemented", cmdToken.val)
	}

	cmd, err := fn(lexer, config)
	if err != nil {
		return nil, err
	}
	command.Executer = cmd

	return command, nil
}

// type Fullscreen struct{}

// func (f Fullscreen) Run(ctx context.Context) error {
// 	return nil
// }

// type Floating string

// func (f Floating) Run(ctx context.Context) error {
// 	return nil
// }

// type Kill struct{}

// func (k Kill) Run(ctx context.Context) error {
// 	return nil
// }

// type Reload struct{}

// func (r Reload) Run(ctx context.Context) error {
// 	return nil
// }

// type WorkspaceCmd struct {
// 	Name   string
// 	Output string
// }

// func (w *WorkspaceCmd) Run(ctx context.Context) error {
// 	return nil
// }

// type MoveDirection string

// func (m MoveDirection) Run(ctx context.Context) error {
// 	return nil
// }

// type MoveContainer struct {
// }

// func (m MoveContainer) Run(ctx context.Context) error {
// 	return nil
// }

// type Split string

// func (s Split) Run(ctx context.Context) error {
// 	return nil
// }

// type Layout string

// func (l Layout) Run(ctx context.Context) error {
// 	return nil
// }

// type Focus string

// func (f Focus) Run(ctx context.Context) error {
// 	return nil
// }

// type FocusOutput string

// func (f FocusOutput) Run(ctx context.Context) error {
// 	return nil
// }

// func Sticky(arg ...string) error {
// 	return nil
// }

// func Rename(arg ...string) error {
// 	return nil
// }

// type Mode string

// func (m Mode) Run(ctx context.Context) error {
// 	return nil
// }
