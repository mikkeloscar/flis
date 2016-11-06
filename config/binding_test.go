package config

import (
	"testing"

	wlc "github.com/mikkeloscar/go-wlc"
	xkb "github.com/mikkeloscar/go-xkbcommon"
)

func TestCmpBinding(t *testing.T) {
	bindings := []struct {
		a   *Binding
		b   *Binding
		ret compare
	}{
		// Shift+A == Shift+A
		{
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			equal,
		},
		// Shift+A+B < Shift+A
		{
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA, xkb.KeyB},
			},
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			less,
		},
		// Shift+A > Shitf+A+B
		{
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA, xkb.KeyB},
			},
			bigger,
		},
		// Shift+A > Shift+Ctrl+A
		{
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			&Binding{
				Modifiers: uint32(wlc.BitModShift) | uint32(wlc.BitModCtrl),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			bigger,
		},
		// Ctrl+A > Shift+A
		{
			&Binding{
				Modifiers: uint32(wlc.BitModCtrl),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			bigger,
		},
		// Shift+A < Ctrl+A
		{
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			&Binding{
				Modifiers: uint32(wlc.BitModCtrl),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			less,
		},
		// Shift+A < Shift+B
		{
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyB},
			},
			less,
		},
		// Shift+B > Shift+A
		{
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyB},
			},
			&Binding{
				Modifiers: uint32(wlc.BitModShift),
				Keys:      []xkb.KeySym{xkb.KeyA},
			},
			bigger,
		},
	}

	for _, x := range bindings {
		ret := cmpBinding(x.a, x.b)
		if ret != x.ret {
			t.Errorf("expected result %d, got %d", ret, x.ret)
		}
	}
}

func TestBindingEqual(t *testing.T) {
	a := &Binding{
		Modifiers: uint32(wlc.BitModShift),
		Keys:      []xkb.KeySym{xkb.KeyG},
	}
	b := &Binding{
		Modifiers: uint32(wlc.BitModShift),
		Keys:      []xkb.KeySym{xkb.KeyG},
	}

	if !a.Equal(b) {
		t.Errorf("expected binding a to be equal to binding b")
	}
}

func TestBindingBigger(t *testing.T) {
	a := &Binding{
		Modifiers: uint32(wlc.BitModShift),
		Keys:      []xkb.KeySym{xkb.KeyG},
	}
	b := &Binding{
		Modifiers: uint32(wlc.BitModShift),
		Keys:      []xkb.KeySym{xkb.KeyG, xkb.KeyH},
	}

	if !a.Bigger(b) {
		t.Errorf("expected binding a to be bigger than binding b")
	}
}

func TestBindingsLen(t *testing.T) {
	bindings := []struct {
		bindings Bindings
		len      int
	}{
		{
			Bindings([]*Binding{}),
			0,
		},
		{
			Bindings([]*Binding{
				{},
			}),
			1,
		},
	}

	for _, b := range bindings {
		if b.len != b.bindings.Len() {
			t.Errorf("expected %d, got %d", b.len, b.bindings.Len())
		}
	}
}

func TestBindingsSwap(t *testing.T) {
	bindings := Bindings([]*Binding{
		{Raw: "1"},
		{Raw: "2"},
	})

	bindings.Swap(0, 1)

	if bindings[0].Raw != "2" && bindings[1].Raw != "1" {
		t.Errorf("expected bindings to be swapped")
	}
}

func TestBindingsLess(t *testing.T) {
	bindings := Bindings([]*Binding{
		{Keys: []xkb.KeySym{xkb.KeyDown, xkb.KeyF}},
		{Keys: []xkb.KeySym{xkb.KeyUp}},
	})

	if !bindings.Less(0, 1) {
		t.Errorf("expected binding 0 to be less than binding 1")
	}
}
