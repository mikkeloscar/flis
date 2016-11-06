package backend

import wlc "github.com/mikkeloscar/go-wlc"

// View is an interface for the wlc View type. Its purpose is to make it easy
// to mock wlc in tests.
type View interface {
	BringAbove(other wlc.View)
	BringToFront()
	Close()
	Focus()
	GetAppID() string
	GetClass() string
	GetGeometry() *wlc.Geometry
	GetMask() uint32
	GetOutput() wlc.Output
	GetParent() wlc.View
	// GetRole() *C.struct_wl_resource
	GetState() uint32
	GetSurface() wlc.Resource
	Title() string
	GetType() uint32
	GetVisibleGeometry() wlc.Geometry
	// GetWlClient() *C.struct_wl_client
	SendBelow(other wlc.View)
	SendToBack()
	SetGeometry(edges uint32, geometry wlc.Geometry)
	SetMask(mask uint32)
	SetOutput(output wlc.Output)
	SetParent(parent wlc.View)
	SetState(state wlc.ViewStateBit, toggle bool)
	SetType(typ wlc.ViewTypeBit, toggle bool)
}
