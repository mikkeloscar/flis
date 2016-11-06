package backend

import wlc "github.com/mikkeloscar/go-wlc"

// Mock implements a mock interface satifying the Backend interface. Can be
// used for tests.
type Mock struct{}

// Exec mocks executing a command in the compositor.
func (m Mock) Exec(bin string, arg ...string) {}

// Terminate mocks terminating the compositor.
func (m Mock) Terminate() {}

// PointerSetPosition mokcs setting pointer position.
func (m Mock) PointerSetPosition(pos wlc.Point) {}
