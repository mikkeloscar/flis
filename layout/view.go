package layout

import (
	"github.com/mikkeloscar/flis/backend"
	wlc "github.com/mikkeloscar/go-wlc"
)

// View defines a view in the layout.
type View struct {
	backend.View
	geometry wlc.Geometry
	parent   Container

	// Name    string
	// Class   string
	// AppID   string
	visible bool
	focused bool
	Border  *Border
}

// NewView initializes a new view container.
func NewView(backend backend.View, parent Container) *View {
	return &View{
		backend,
		wlc.GeometryZero, // TODO: set default
		parent,
		false,
		false,
		nil, // TODO: init border
	}
}

// Type returns the view container type.
func (v *View) Type() ContainerType {
	return CView
}

// Geometry returns the view geometry.
// TODO: with or without border?
func (v *View) Geometry() *wlc.Geometry {
	return &v.geometry
}

// Children returns nil since a view doesn't have any children.
func (v *View) Children() []Container {
	return nil
}

// Floating returns nil since a view doesn't have any floating children.
func (v *View) Floating() []Container {
	return nil
}

// Focused returns a pointer to itself if it's focused.
func (v *View) Focused() Container {
	if v.focused {
		return v
	}
	return nil
}

// Fullscreen returns true if view is in fullscreen mode.
func (v *View) Fullscreen() bool {
	return v.GetState() == wlc.BitFullscreen
}

// Parent returns the parent container of the view.
func (v *View) Parent() Container {
	return v.parent
}

// Visible returns true if the view is visible. The parent container must be
// visible for the view to be visible.
func (v *View) Visible() bool {
	return v.parent.Visible() && v.visible
}

// SetVisible sets the visibility state of the view.
func (v *View) SetVisible(visible bool) {
	v.visible = visible
}

// AddChild is a no-op function since views can't have any children.
func (v *View) AddChild(_ Container) {}

// Border defines a view border.
// TODO: implement.
type Border struct {
	Geometry wlc.Geometry
}
