package input

import (
	"testing"

	xkb "github.com/mikkeloscar/go-xkbcommon"
)

// TestIsPressed tests the IsPressed method.
func TestIsPressed(t *testing.T) {
	state := NewKeyState()

	state.keySyms[1] = struct{}{}
	state.keyCodes[2] = struct{}{}

	if !state.IsPressed(1, 0) {
		t.Errorf("should be pressed indentified by keysym")
	}

	if !state.IsPressed(0, 2) {
		t.Errorf("should be pressed indentified by keycode")
	}

	if state.IsPressed(0, 0) {
		t.Errorf("should not be pressed")
	}
}

// TestPressKey tests if a key gets pressed.
func TestPressKey(t *testing.T) {
	state := NewKeyState()

	state.PressKey(xkb.Keyspace, 1)
	state.PressKey(0, 3)
	state.PressKey(xkb.KeyTab, 0)

	if !state.IsPressed(xkb.Keyspace, 0) {
		t.Errorf("space keysym should be pressed")
	}

	if !state.IsPressed(0, 3) {
		t.Errorf("keycode should be pressed")
	}

	if state.IsPressed(xkb.KeyTab, 0) {
		t.Errorf("tab should not be pressed")
	}
}

// TestReleaseKey tests if a key gets released.
func TestReleaseKey(t *testing.T) {
	state := NewKeyState()
	state.PressKey(xkb.Keyspace, 1)
	if !state.IsPressed(xkb.Keyspace, 0) {
		t.Errorf("space keysym should be pressed")
	}

	state.ReleaseKey(xkb.Keyspace, 1)
	if state.IsPressed(xkb.Keyspace, 0) {
		t.Errorf("space keysym should not be pressed")
	}
}
