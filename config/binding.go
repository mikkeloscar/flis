package config

import (
	xkb "github.com/mikkeloscar/go-xkbcommon"
)

type compare int

const (
	equal compare = iota
	less
	bigger
)

// Binding defines a keybinding with related command.
type Binding struct {
	Modifiers uint32
	Keys      []xkb.KeySym
	Command   *Command
	Raw       string
}

// Equal returns true if two keybindings are identical.
func (b *Binding) Equal(a *Binding) bool {
	return cmpBinding(b, a) == equal
}

// Less returns true if binding b is less than binding a.
func (b *Binding) Less(a *Binding) bool {
	return cmpBinding(b, a) == less
}

// Bigger returns true if binding b is bigger than binding a.
func (b *Binding) Bigger(a *Binding) bool {
	return cmpBinding(b, a) == bigger
}

// cmpBinding compares two bindings a and b. A binding is considered less than
// the other if it has more keys+modifiers defined. If both bindings
// have the same number of keys and modifiers defined but not the same
// keys/modifiers then the keys with the highest value will be considered the
// biggest binding.
func cmpBinding(a, b *Binding) compare {
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

	if (len(b.Keys) + modB) > (len(a.Keys) + modA) {
		return bigger
	}
	if (len(b.Keys) + modB) < (len(a.Keys) + modA) {
		return less
	}

	if a.Modifiers > b.Modifiers {
		return bigger
	} else if a.Modifiers < b.Modifiers {
		return less
	}

	for i, keyA := range a.Keys {
		if keyA > b.Keys[i] {
			return bigger
		}

		if keyA < b.Keys[i] {
			return less
		}
	}

	return equal
}

// Bindings is a list of keybindings.
type Bindings []*Binding

// Len returns the length of the binding list.
func (b Bindings) Len() int {
	return len(b)
}

// Swap swaps two bindings in the bindings list.
func (b Bindings) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Less returns true if binding at index i should be sorted before binding at
// index j.
func (b Bindings) Less(i, j int) bool {
	return b[i].Less(b[j])
}
