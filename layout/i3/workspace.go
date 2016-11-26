package i3

import (
	"context"

	"github.com/mikkeloscar/flis/layout"
)

// NewWorkspace creates a new workspace for an output.
func (l *Layout) NewWorkspace(ctx context.Context, output *layout.Output, name string, num uint) {
	workspace := layout.NewWorkspace(name, num, output)
	output.AddChild(workspace)
}
