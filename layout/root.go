package layout

import (
	"reflect"
	"sort"

	log "github.com/sirupsen/logrus"
	wlc "github.com/mikkeloscar/go-wlc"
)

// Root is the root container in a layout. The root container manages the
// output containers.
type Root struct {
	outputs map[string]*Output
	focused *Output
}

// NewRoot initializes a new empty root container.
func NewRoot() *Root {
	return &Root{
		outputs: make(map[string]*Output),
		focused: nil,
	}
}

// Type returns the root container type.
func (r *Root) Type() ContainerType {
	return CRoot
}

// Geometry returns the root container geometry which is always nil.
func (r *Root) Geometry() *wlc.Geometry {
	return nil
}

// Children returns a list for output containers attached to the root
// container.
func (r *Root) Children() []Container {
	containers := make([]Container, 0, len(r.outputs))
	for _, o := range r.outputs {
		containers = append(containers, o)
	}
	return containers
}

// Floating always returns nil because the root container can't have floating
// child containers.
func (r *Root) Floating() []Container {
	return nil
}

// Focused returns the focused output container.
func (r *Root) Focused() Container {
	return r.focused
}

// Parent returns nil because the root container has no parents.
func (r *Root) Parent() Container {
	return nil
}

// Visible returns true if the container is visible. The root container is
// always visible.
func (r *Root) Visible() bool {
	return true
}

// AddChild adds an output to the root container.
func (r *Root) AddChild(output Container) {
	switch o := output.(type) {
	case *Output:
		log.Debugf("Added output %s", o.Name())
		r.outputs[o.Name()] = o

		// if there is only one output then focus it.
		if len(r.outputs) == 1 {
			r.focused = o
		}
	default:
		log.Errorf("Failed to add output, invalid container type: %s", reflect.TypeOf(output))
	}
}

// SortedWorkspaces returns an aggregated sorted list of workspaces on all
// outputs.
func (r *Root) SortedWorkspaces() []*Workspace {
	var ws []*Workspace
	for _, output := range r.outputs {
		ws = append(ws, output.workspaces...)
	}

	sort.Sort(Workspaces(ws))
	return ws
}
