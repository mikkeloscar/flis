package compositor

import (
	"testing"

	"github.com/mikkeloscar/flise/backend"
	"github.com/mikkeloscar/flise/config"
	wlc "github.com/mikkeloscar/go-wlc"
	xkb "github.com/mikkeloscar/go-xkbcommon"
)

// Test PointerMotion cb.
func TestPointerMotion(t *testing.T) {
	c := New(nil, backend.Mock{}, nil)
	if !c.PointerMotion(0, 0, &wlc.PointZero) {
		t.Errorf("expected pointer motion to succeed")
	}
}

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
