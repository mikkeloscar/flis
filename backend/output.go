package backend

import wlc "github.com/mikkeloscar/go-wlc"

// Output is an interface for the wlc Output type. Its purpose is to make it
// easy to mock wlc in tests.
type Output interface {
	Focus()
	GetMask() uint32
	GetMutableViews() []wlc.View // TODO: abstract this further if needed.
	Name() string
	GetRenderer() wlc.Renderer
	GetResolution() *wlc.Size
	GetVirtualResolution() *wlc.Size
	SetResolution(resolution wlc.Size, scale uint32)
	GetScale() uint32
	GetSleep() bool
	GetViews() []wlc.View
	ScheduleRender()
	SetMask(mask uint32)
	SetSleep(sleep bool)
	SetViews(views []wlc.View) bool
}
