package backend

import (
	"testing"

	wlc "github.com/mikkeloscar/go-wlc"
)

func TestMock(t *testing.T) {
	m := Mock{}

	m.Exec("command")
	m.Terminate()
	m.PointerSetPosition(wlc.PointZero)
}
