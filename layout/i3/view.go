package i3

import (
	"github.com/mikkeloscar/flise/backend"
	"github.com/mikkeloscar/flise/context"
	"github.com/mikkeloscar/flise/layout"
)

func (l *Layout) NewView(ctx context.Context, view backend.View) {
	// TODO: wlc, set size
	// TODO: check if sibling or parent
	parent := l.FocusedByType(ctx, layout.CWorkspace)
	v := layout.NewView(view, parent)
	parent.AddChild(v)
}
