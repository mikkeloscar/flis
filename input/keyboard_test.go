package input

import (
	"testing"

	xkb "github.com/mikkeloscar/go-xkbcommon"
	"github.com/stretchr/testify/assert"
)

// TestIsPressed tests the IsPressed method.
func TestIsPressed(t *testing.T) {
	state := NewKeyState()

	state.keySyms[1] = struct{}{}
	state.keyCodes[2] = struct{}{}

	assert.True(t, state.IsPressed(1, 0), "Should be pressed indentified by keysym")
	assert.True(t, state.IsPressed(0, 2), "Should be pressed indentified by keycode")
	assert.False(t, state.IsPressed(0, 0), "Should not be pressed")
}

// TestPressKey tests if a key gets pressed.
func TestPressKey(t *testing.T) {
	state := NewKeyState()

	state.PressKey(xkb.Keyspace, 1)
	state.PressKey(0, 3)
	state.PressKey(xkb.KeyTab, 0)

	assert.True(t, state.IsPressed(xkb.Keyspace, 0), "Space keysym should be pressed")
	assert.True(t, state.IsPressed(0, 3), "Keycode should be pressed")
	assert.False(t, state.IsPressed(xkb.KeyTab, 0), "Tab should not be pressed")
}

// TestReleaseKey tests if a key gets released.
func TestReleaseKey(t *testing.T) {
	state := NewKeyState()
	state.PressKey(xkb.Keyspace, 1)
	assert.True(t, state.IsPressed(xkb.Keyspace, 0), "Space keysym should be pressed")

	state.ReleaseKey(xkb.Keyspace, 1)
	assert.False(t, state.IsPressed(xkb.Keyspace, 0), "Space keysym should not be pressed")
}
