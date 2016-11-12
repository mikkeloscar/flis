package layout

import (
	"testing"

	wlc "github.com/mikkeloscar/go-wlc"
)

// mockView mocks the backend.View interface.
type mockView struct{}

func (m mockView) BringAbove(other wlc.View)                       {}
func (m mockView) BringToFront()                                   {}
func (m mockView) Close()                                          {}
func (m mockView) Focus()                                          {}
func (m mockView) GetAppID() string                                { return "" }
func (m mockView) GetClass() string                                { return "" }
func (m mockView) GetGeometry() *wlc.Geometry                      { return nil }
func (m mockView) GetMask() uint32                                 { return 0 }
func (m mockView) GetOutput() wlc.Output                           { return 0 }
func (m mockView) GetParent() wlc.View                             { return 0 }
func (m mockView) GetState() uint32                                { return 0 }
func (m mockView) GetSurface() wlc.Resource                        { return 0 }
func (m mockView) Title() string                                   { return "" }
func (m mockView) GetType() uint32                                 { return 0 }
func (m mockView) GetVisibleGeometry() wlc.Geometry                { return wlc.GeometryZero }
func (m mockView) SendBelow(other wlc.View)                        {}
func (m mockView) SendToBack()                                     {}
func (m mockView) SetGeometry(edges uint32, geometry wlc.Geometry) {}
func (m mockView) SetMask(mask uint32)                             {}
func (m mockView) SetOutput(output wlc.Output)                     {}
func (m mockView) SetParent(parent wlc.View)                       {}
func (m mockView) SetState(state wlc.ViewStateBit, toggle bool)    {}
func (m mockView) SetType(typ wlc.ViewTypeBit, toggle bool)        {}

// mockViewFullscreen mocks the backend.View interface with state BitFullscreen.
type mockViewFullscreen struct {
	mockView
}

func (m mockViewFullscreen) GetState() uint32 {
	return wlc.BitFullscreen
}

// TestViewType tests getting the view type from a view container.
func TestViewType(t *testing.T) {
	v := NewView(wlc.View(0), nil)
	if v.Type() != CView {
		t.Errorf("expected container type CView, got %d", v.Type())
	}
}

// TestViewGeometry tests getting geometry of view.
func TestViewGeometry(t *testing.T) {
	v := NewView(wlc.View(0), nil)
	if !wlc.GeometryEquals(*v.Geometry(), wlc.GeometryZero) {
		t.Errorf("expected geometry (0,0 - 0x0), got %v", v.Geometry())
	}
}

// TestViewChildren tests getting 0 children from a view container.
func TestViewChildren(t *testing.T) {
	v := NewView(wlc.View(0), nil)
	if len(v.Children()) != 0 {
		t.Errorf("expected 0 child containers, got %d", len(v.Children()))
	}
}

// TestViewFloating tests getting 0 floating child containers from view container.
func TestViewFloating(t *testing.T) {
	v := NewView(wlc.View(0), nil)
	if len(v.Floating()) != 0 {
		t.Errorf("expected 0 floating containers, got %d", len(v.Floating()))
	}
}

// TestViewFocused tests getting focused container.
func TestViewFocused(t *testing.T) {
	v := NewView(wlc.View(0), nil)
	if v.Focused() != nil {
		t.Errorf("did not expect to get a focused container, got %v", v.Focused())
	}

	v.focused = true
	if v.Focused() != v {
		t.Errorf("expected to get container %v, got %v", v, v.Focused())
	}
}

// TestViewFullscreen tests if view is fullscreen.
func TestViewFullscreen(t *testing.T) {
	v := NewView(mockViewFullscreen{}, nil)
	if !v.Fullscreen() {
		t.Errorf("expected view to be fullscreen")
	}

	v = NewView(mockView{}, nil)
	if v.Fullscreen() {
		t.Errorf("did not expect view to be fullscreen")
	}
}

// TestViewParent tests getting parent container of view.
func TestViewParent(t *testing.T) {
	w := NewWorkspace("1", 1, nil)
	v := NewView(wlc.View(0), w)
	w.AddChild(v)
	if v.Parent() != w {
		t.Errorf("expected to get parent %v, got %v", w, v.Parent())
	}
}

// TestViewVisible tests if workspace is visible.
func TestViewVisible(t *testing.T) {
	r := NewRoot()
	o := NewOutput(wlc.Output(0), r)
	r.AddChild(o)
	w := NewWorkspace("1", 1, o)
	o.AddChild(w)
	v := NewView(wlc.View(0), w)
	w.AddChild(v)
	if v.Visible() {
		t.Errorf("did not expect view to be visible")
	}

	v.SetVisible(true)
	if !v.Visible() {
		t.Errorf("expected view to be visible")
	}

	v.SetVisible(false)
	if v.Visible() {
		t.Errorf("did not expect view to be visible")
	}
}

// TestViewAddChild tests AddChild no-op.
func TestViewAddChild(t *testing.T) {
	v := NewView(wlc.View(0), nil)
	v.AddChild(nil)
}
