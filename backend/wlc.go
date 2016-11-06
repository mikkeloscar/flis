package backend

import wlc "github.com/mikkeloscar/go-wlc"

// WLC is a wrapper around the wlc api which implement the Backend interface.
type WLC struct{}

// Exec executes a command in the compositor.
func (w WLC) Exec(bin string, arg ...string) {
	wlc.Exec(bin, arg...)
}

// Terminate terminates the compositor.
func (w WLC) Terminate() {
	wlc.Terminate()
}

// PointerSetPosition sets pointer position.
func (w WLC) PointerSetPosition(pos wlc.Point) {
	wlc.PointerSetPosition(pos)
}
