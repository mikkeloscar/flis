package commands

import (
	"fmt"

	"github.com/mikkeloscar/flise/context"
)

// command is an interface for commands.
type command interface {
	Exec(context.Context) error
	String() string
}

// Criteria - TODO: not implemented.
type Criteria int

// Command defines a command with criteria and possible chained sub commands.
type Command struct {
	command
	// TODO: criteria
	Criteria *Criteria
	Next     *Command
}

// String returns the string representation of a command.
func (c *Command) String() string {
	// TODO: criteria
	if c.Next != nil {
		return c.command.String() + ", " + c.Next.String()
	}
	return c.command.String()
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

type commandParser func(lex *lexer) (command, error)

// table mapping command names to parse functions able to parse the command
// definitions.
var cmdParseTable = map[string]commandParser{
	"exec": parseExec,
	"exit": parseExit,
}

// Parse parses a command string into a command structure.
func Parse(commandStr string) (*Command, error) {
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

	cmd, err := fn(lexer)
	if err != nil {
		return nil, err
	}
	command.command = cmd

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
