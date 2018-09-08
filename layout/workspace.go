package layout

import (
	"reflect"

	log "github.com/sirupsen/logrus"
	wlc "github.com/mikkeloscar/go-wlc"
)

// Workspace is a workspace container in the layout. A workspace has a list of
// tiled containers plus a list of floating containers associated with it.
type Workspace struct {
	name       string
	Num        uint
	containers []Container
	floating   []Container
	focused    Container
	output     *Output
}

// NewWorkspace initializes a new empty workspace container.
func NewWorkspace(name string, num uint, output *Output) *Workspace {
	return &Workspace{
		name:       name,
		Num:        num,
		containers: make([]Container, 0),
		floating:   make([]Container, 0),
		focused:    nil,
		output:     output,
	}
}

// Type returns the workspace container type.
func (w *Workspace) Type() ContainerType {
	return CWorkspace
}

// Name return the name of the workspace.
func (w *Workspace) Name() string {
	return w.name
}

// Geometry for the workspace is the geometry of the parent output.
func (w *Workspace) Geometry() *wlc.Geometry {
	return w.output.Geometry()
}

// Children returns a list of (non-floating) containers on the workspace.
func (w *Workspace) Children() []Container {
	return w.containers
}

// Floating returns a list of floating containers on the workspace.
func (w *Workspace) Floating() []Container {
	return w.floating
}

// Focused returns the focused child container of the workspace.
func (w *Workspace) Focused() Container {
	return w.focused
}

// Parent returns the parent output of the workspace.
func (w *Workspace) Parent() Container {
	return w.output
}

// AddChild adds a child container to the workspace.
func (w *Workspace) AddChild(container Container) {
	switch c := container.(type) {
	case *View:
		w.containers = append(w.containers, c)

		// focus added container
		w.focused = c
		log.Debugf("Added container '%s' to workspace %s", c.Title(), w.Name())
	default:
		log.Errorf("Failed to add container, invalid container type: %s", reflect.TypeOf(container))
	}
}

// Visible returns true if workspace is visible.
func (w *Workspace) Visible() bool {
	return w.output.Visible() && w.output.focused == w
}

// Workspaces is a list of workspaces.
type Workspaces []*Workspace

// Len returns the length of the workspace list.
func (w Workspaces) Len() int {
	return len(w)
}

// Swap swaps two workspaces in the workspace list.
func (w Workspaces) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

// Less returns true if workspace at index i should be sorted before workspace
// at index j.
func (w Workspaces) Less(i, j int) bool {
	return w[i].Num < w[j].Num
}
