package i3

import (
	"github.com/mikkeloscar/flis/backend"
	"github.com/mikkeloscar/flis/context"
	"github.com/mikkeloscar/flis/layout"
)

// NewOutput creates a new output from a wlc output handle.
func (l *Layout) NewOutput(ctx context.Context, output backend.Output) {
	// size = output.GetResolution() // TODO: wlc, get size
	// TODO: create+assing workspace
	o := layout.NewOutput(output, l.root)
	l.root.AddChild(o)
}
