package backend

import (
	"github.com/mikkeloscar/flise/context"
	wlc "github.com/mikkeloscar/go-wlc"
)

// Backend is an inteface describing the functions of a backend like wlc. This
// is used instead of calling wlc functions directly making it easier to
// implement tests which mocks the functionality of wlc.
type Backend interface {
	Exec(bin string, arg ...string)
	Terminate()
	PointerSetPosition(pos wlc.Point)
	KeyboardGetKeysymForKey(key uint32, mods *wlc.Modifiers) uint32
}

// Get backend from context.
func Get(ctx context.Context) Backend {
	return ctx.MustGet("backend").(Backend)
}
