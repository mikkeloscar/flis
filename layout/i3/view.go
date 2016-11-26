package i3

import (
	"context"

	"github.com/mikkeloscar/flis/backend"
	"github.com/mikkeloscar/flis/layout"
)

// NewView adds a new view to the layout. The view will be added to the
// currently focused container.
func (l *Layout) NewView(ctx context.Context, view backend.View) {
	// TODO: check if sibling or parent
	parent := l.FocusedByType(ctx, layout.CWorkspace)
	v := layout.NewView(view, parent)
	parent.AddChild(v)
}
