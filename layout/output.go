package layout

import (
	"reflect"
	"sort"

	log "github.com/Sirupsen/logrus"
	"github.com/mikkeloscar/flis/backend"
	wlc "github.com/mikkeloscar/go-wlc"
)

// Output describes an output in the layout.
type Output struct {
	backend.Output
	workspaces []*Workspace
	focused    *Workspace
	root       *Root
}

// NewOutput initializes a new empty output.
func NewOutput(output backend.Output, root *Root) *Output {
	return &Output{
		output,
		make([]*Workspace, 0),
		nil,
		root,
	}
}

// Type returns the output container type.
func (o *Output) Type() ContainerType {
	return COutput
}

// Geometry returns the geometry of the output.
func (o *Output) Geometry() *wlc.Geometry {
	// TODO: relative to other outputs
	return &wlc.Geometry{
		Origin: wlc.PointZero,
		Size:   *o.GetVirtualResolution(),
	}
}

// func (o *Output) Resolution() *wlc.Size {
// 	return o.GetVirtualResolution()
// }

func (o *Output) Children() []Container {
	containers := make([]Container, len(o.workspaces))
	for i, w := range o.workspaces {
		containers[i] = w
	}
	return containers
}

func (o *Output) Floating() []Container {
	return nil
}

func (o *Output) Focused() Container {
	return o.focused
}

func (o *Output) Parent() Container {
	return o.root
}

func (o *Output) AddChild(workspace Container) {
	switch w := workspace.(type) {
	case *Workspace:
		o.workspaces = append(o.workspaces, w)
		// sort workspaces by num
		sort.Sort(Workspaces(o.workspaces))

		// if there is only one workspace on the output then focus it.
		if len(o.workspaces) == 1 {
			o.focused = w
		}
		log.Debugf("Added workspace '%s' to output %s", w.Name(), o.Name())
	default:
		log.Errorf("Failed to add workspace, invalid container type: %s", reflect.TypeOf(workspace))
	}
}

func (o *Output) Visible() bool {
	// TODO: handle possible invisible cases
	return true
}
