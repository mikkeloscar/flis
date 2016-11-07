package compositor

import (
	"fmt"
	"testing"

	"github.com/mikkeloscar/flis/backend"
	"github.com/mikkeloscar/flis/config"
	"github.com/mikkeloscar/flis/config/commands"
	"github.com/mikkeloscar/flis/context"
	"github.com/mikkeloscar/flis/layout"
	wlc "github.com/mikkeloscar/go-wlc"
	xkb "github.com/mikkeloscar/go-xkbcommon"
)

// Test OutputCreated cb.
func TestOutputCreated(t *testing.T) {
	c := New(nil, nil, layout.Mock{})
	if !c.OutputCreated(0) {
		t.Errorf("expected output created to succeed")
	}
}

// Test ViewCreated cb.
func TestViewCreated(t *testing.T) {
	c := New(nil, nil, layout.Mock{})
	if !c.ViewCreated(0) {
		t.Errorf("expected view created to succeed")
	}
}

// Test PointerMotion cb.
func TestPointerMotion(t *testing.T) {
	c := New(nil, backend.Mock{}, nil)
	if !c.PointerMotion(0, 0, &wlc.PointZero) {
		t.Errorf("expected pointer motion to succeed")
	}
}

type mockExecuter struct{}

func (m mockExecuter) Exec(ctx context.Context) error {
	return nil
}

func (m mockExecuter) String() string {
	return "mock"
}

type mockExecuterFail struct {
	mockExecuter
}

func (m mockExecuterFail) Exec(ctx context.Context) error {
	return fmt.Errorf("failed to execute")
}

// Test KeyboardKey cb.
func TestKeyboardKey(t *testing.T) {
	conf := config.New()
	binding := &config.Binding{
		Keys: []xkb.KeySym{xkb.KeyA},
		Command: &commands.Command{
			Executer: mockExecuter{},
		},
	}
	conf.AddBinding("default", binding)
	c := New(conf, backend.Mock{}, nil)

	modifiers := wlc.Modifiers{
		Leds: 0x1,
		Mods: 0x0,
	}

	if ret := c.KeyboardKey(0, 0, modifiers, xkb.KeyA, wlc.KeyStatePressed); ret != eventHandled {
		t.Errorf("expected eventHandled(%t) return value, got %t", eventHandled, ret)
	}

	// set failing binding command executer.
	binding.Command.Executer = mockExecuterFail{}
	conf.AddBinding("default", binding)

	if ret := c.KeyboardKey(0, 0, modifiers, xkb.KeyA, wlc.KeyStatePressed); ret != eventHandled {
		t.Errorf("expected eventHandled(%t) return value, got %t", eventHandled, ret)
	}

	if ret := c.KeyboardKey(0, 0, modifiers, xkb.KeyA, wlc.KeyStateReleased); ret != eventPassthrough {
		t.Errorf("expected eventPassthrough(%t) return value, got %t", eventPassthrough, ret)
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
