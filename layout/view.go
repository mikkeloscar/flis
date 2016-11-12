package layout

import (
	"github.com/mikkeloscar/flis/backend"
	wlc "github.com/mikkeloscar/go-wlc"
)

type View struct {
	backend.View
	geometry wlc.Geometry
	parent   Container

	// Name    string
	// Class   string
	// AppID   string
	visible bool
	Border  *Border
}

func NewView(backend backend.View, parent Container) *View {
	return &View{
		backend,
		wlc.GeometryZero, // TODO: set default
		parent,
		false,
		nil, // TODO: init border
	}
}

func (v *View) Type() ContainerType {
	return CView
}

func (v *View) Geometry() *wlc.Geometry {
	return &v.geometry
}

func (v *View) Children() []Container {
	return nil
}

func (v *View) Floating() []Container {
	return nil
}

func (v *View) Focused() Container {
	return nil
}

// Fullscreen returns true if view is in fullscreen mode.
func (v *View) Fullscreen() bool {
	return v.GetState() == wlc.BitFullscreen
}

func (v *View) Parent() Container {
	return v.parent
}

func (v *View) Visible() bool {
	return v.visible
}

func (v *View) SetVisible(visible bool) {
	v.visible = visible
}

func (v *View) AddChild(_ Container) {}

// TODO: implement
type Border struct {
	Geometry wlc.Geometry
}
