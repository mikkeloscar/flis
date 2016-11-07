package i3

import (
	"github.com/mikkeloscar/flis/context"
	"github.com/mikkeloscar/flis/layout"
)

// NewWorkspace creates a new workspace for an output.
func (l *Layout) NewWorkspace(ctx context.Context, output *layout.Output, name string) {
	workspace := layout.NewWorkspace(name, output)
	output.AddChild(workspace)
}
