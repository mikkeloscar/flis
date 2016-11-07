package backend

import (
	"testing"

	wlc "github.com/mikkeloscar/go-wlc"
)

// TestMock creates test coverage for the backend mock.
func TestMock(t *testing.T) {
	m := Mock{}

	m.Exec("command")
	m.Terminate()
	m.PointerSetPosition(wlc.PointZero)
	m.KeyboardGetKeysymForKey(0, nil)
}
