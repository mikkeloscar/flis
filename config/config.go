package config

import (
	"sort"

	"github.com/mikkeloscar/go-wlc"
	"github.com/mikkeloscar/go-xkbcommon"
)

type Config struct {
	vars       map[string]string
	mode       string
	Modes      map[string][]*Binding
	Workspaces []Workspace
	Bars       []Bar
	Outputs    map[string]Output
}

func (c *Config) Bindings() []*Binding {
	if b, ok := c.Modes[c.mode]; ok {
		return b
	}

	return nil
}

func (c *Config) AddBinding(b *Binding) bool {
	var bindings []*Binding
	if bs, ok := c.Modes[c.mode]; ok {
		bindings = bs
	} else {
		c.Modes[c.mode] = []*Binding{b}
		sort.Sort(Bindings(c.Modes[c.mode]))
		return false
	}

	for i, binding := range bindings {
		if b.Equal(binding) {
			bindings = append(bindings[:i-1], bindings[i+1:]...)
			sort.Sort(Bindings(bindings))
			return true
		}
	}

	return false
}

type Bindings []*Binding

func (b Bindings) Len() int {
	return len(b)
}

func (b Bindings) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Bindings) Less(i, j int) bool {
	return b[i].Less(b[j])
}

type Binding struct {
	Modifiers uint32
	Keys      []xkb.KeySym
	Command   *Command
}

type Criteria int

func (b *Binding) Equal(a *Binding) bool {
	return cmpBinding(b, a) == 0
}

func (b *Binding) Less(a *Binding) bool {
	return cmpBinding(b, a) == -1
}

func (b *Binding) Bigger(a *Binding) bool {
	return cmpBinding(b, a) == 1
}

func cmpBinding(a, b *Binding) int {
	modA := 0
	modB := 0

	for i := 0; i < 8; i++ {
		if a.Modifiers&(1<<uint(i)) != 0 {
			modA++
		}

		if b.Modifiers&(1<<uint(i)) != 0 {
			modB++
		}
	}

	if (len(b.Keys) + modB) != (len(a.Keys) + modA) {
		return (len(b.Keys) + modB) - (len(a.Keys) + modA)
	}

	if a.Modifiers > b.Modifiers {
		return 1
	} else if a.Modifiers < b.Modifiers {
		return -1
	}

	for i, keyA := range a.Keys {
		if keyA > b.Keys[i] {
			return 1
		}

		if keyA < b.Keys[i] {
			return -1
		}
	}

	return 0
}

type Workspace struct {
	Number uint
	Name   string
	Active bool
}

type Bar struct {
	ID string
}

type Command struct {
	Args     []string
	Criteria *Criteria
	Exec     func(arg ...string) error
	Next     *Command
}

func (c *Command) Run() error {
	err := c.Exec(c.Args...)
	if err != nil {
		return err
	}

	if c.Next != nil {
		return c.Next.Run()
	}

	return nil
}

type Output struct {
	Name       string
	Enabled    bool
	Size       wlc.Size
	Pos        wlc.Point
	Background string
}
