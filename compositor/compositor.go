package compositor

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/mikkeloscar/flis/backend"
	"github.com/mikkeloscar/flis/config"
	"github.com/mikkeloscar/flis/input"
	"github.com/mikkeloscar/flis/layout"
	wlc "github.com/mikkeloscar/go-wlc"
	xkb "github.com/mikkeloscar/go-xkbcommon"
)

const (
	eventPassthrough = false
	eventHandled     = true
)

// Compositor defines the state of the running compositor.
type Compositor struct {
	layout   layout.Layout
	ctx      context.Context
	keyState *input.KeyState
	backend  backend.Backend
}

// New initializes a new Compositor.
func New(conf *config.Config, backend backend.Backend, layout layout.Layout) *Compositor {
	ctx := context.WithValue(context.Background(), "config", conf)
	ctx = context.WithValue(ctx, "backend", backend)
	ctx = context.WithValue(ctx, "layout", layout)

	return &Compositor{
		layout:   layout,
		keyState: input.NewKeyState(),
		ctx:      ctx,
		backend:  backend,
	}
}

// OutputCreated is the callback triggered when an output is added by the
// backend.
func (c *Compositor) OutputCreated(o wlc.Output) bool {
	c.layout.NewOutput(c.ctx, o)

	// TODO: return false on failure
	return true
}

// ViewCreated is the callback triggered when a view is added by the backend.
func (c *Compositor) ViewCreated(v wlc.View) bool {
	c.layout.NewView(c.ctx, v)

	// TODO: return false on failure
	return true
}

// PointerMotion is the callback triggered when the pointer is moved.
func (c *Compositor) PointerMotion(view wlc.View, time uint32, pos *wlc.Point) bool {
	c.backend.PointerSetPosition(*pos)
	return true
}

// KeyboardKey is the callback triggered when a key is pressed or released in
// the compositor.
func (c *Compositor) KeyboardKey(view wlc.View, time uint32, modifiers wlc.Modifiers, key uint32, state wlc.KeyState) bool {
	sym := xkb.KeySym(c.backend.KeyboardGetKeysymForKey(key, nil))

	if state == wlc.KeyStatePressed {
		c.keyState.PressKey(sym, key)
	} else {
		c.keyState.ReleaseKey(sym, key)
	}

	conf := config.Get(c.ctx)

	// TODO: maybe move to separate function
	for _, binding := range conf.Bindings() {
		if modifiers.Mods^binding.Modifiers == 0 {
			if c.validBinding(binding) {
				err := binding.Command.Run(c.ctx)
				if err != nil {
					log.Errorf("Failed to run command %s: %s", binding.Command, err)
				}
				return eventHandled
			}
		}
	}

	return eventPassthrough
}

// validBinding returns true if the binding is valid for the current compositor
// keystate.
func (c *Compositor) validBinding(binding *config.Binding) bool {
	for _, key := range binding.Keys {
		if !c.keyState.IsPressed(key, 0) {
			return false
		}
	}

	return true
}
