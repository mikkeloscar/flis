package compositor

import (
	"testing"

	"github.com/mikkeloscar/flise/config"
	xkb "github.com/mikkeloscar/go-xkbcommon"
)

// Test validBinding function.
func TestValidBinding(t *testing.T) {
	c := New(nil, nil, nil)
	binding := &config.Binding{
		Keys: []xkb.KeySym{xkb.KeyA},
	}

	c.keyState.PressKey(xkb.KeyA, 1)
	if !c.validBinding(binding) {
		t.Errorf("expected binding to be valid")
	}

	c.keyState.ReleaseKey(xkb.KeyA, 0)
	if c.validBinding(binding) {
		t.Errorf("expected binding to not be valid")
	}
}
