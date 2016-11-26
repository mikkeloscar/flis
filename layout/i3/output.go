package i3

import (
	"context"
	"fmt"

	"github.com/mikkeloscar/flis/backend"
	"github.com/mikkeloscar/flis/config"
	"github.com/mikkeloscar/flis/layout"
)

// NewOutput creates a new output from a backend output handle.
func (l *Layout) NewOutput(ctx context.Context, output backend.Output) {
	// size = output.GetResolution() // TODO: wlc, get size
	config := config.Get(ctx)
	name, num := findNextWorkspace(l.root.SortedWorkspaces(), config.Workspaces)
	o := layout.NewOutput(output, l.root)
	l.root.AddChild(o)
	l.NewWorkspace(ctx, o, name, num)
}

// findNextWorkspace finds next available workspace name & number.
func findNextWorkspace(ws []*layout.Workspace, confWs map[uint]string) (string, uint) {
	num := findAvailableWorkspaceNum(ws)
	if name, ok := confWs[num]; ok {
		return name, num
	}

	return fmt.Sprintf("%d", num), num
}

// findAvailableWorkspaceNum finds the first available (unused) workspace
// number based on the list of currently allocated workspaces.
func findAvailableWorkspaceNum(ws []*layout.Workspace) uint {
	var num uint
	for i := 0; i < len(ws); i++ {
		num = uint(i + 1)
		if num < ws[i].Num {
			return num
		}
	}
	return num + 1
}
