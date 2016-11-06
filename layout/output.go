package layout

import (
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/mikkeloscar/flise/backend"
	wlc "github.com/mikkeloscar/go-wlc"
)

type Output struct {
	backend.Output
	workspaces []*Workspace
	focused    *Workspace
	root       *Root
}

func NewOutput(output backend.Output, root *Root) *Output {
	return &Output{
		output,
		make([]*Workspace, 0),
		nil,
		root,
	}
}

func (o *Output) Type() ContainerType {
	return COutput
}

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

func (o *Output) Fullscreen() Container {
	if o.focused != nil {
		return o.focused.Fullscreen()
	}
	return nil
}

func (o *Output) Parent() Container {
	return o.root
}

func (o *Output) AddChild(workspace Container) {
	switch w := workspace.(type) {
	case *Workspace:
		// TODO: sort workspaces
		log.Debugf("Added workspace %s for output %d", w.Name(), o.Name())
		o.workspaces = append(o.workspaces, w)
	default:
		log.Errorf("Failed to add workspace, invalid container type: %s", reflect.TypeOf(workspace))
	}
}

func (o *Output) Visible() bool {
	// TODO: handle possible invisible cases
	return true
}
