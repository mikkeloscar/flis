package layout

import (
	"testing"

	wlc "github.com/mikkeloscar/go-wlc"
)

// mockOutput mocks the backend.Output interface.
type mockOutput struct{}

func (m mockOutput) Focus()                                          {}
func (m mockOutput) GetMask() uint32                                 { return 0 }
func (m mockOutput) GetMutableViews() []wlc.View                     { return nil }
func (m mockOutput) Name() string                                    { return "output" }
func (m mockOutput) GetRenderer() wlc.Renderer                       { return 0 }
func (m mockOutput) GetResolution() *wlc.Size                        { return nil }
func (m mockOutput) GetVirtualResolution() *wlc.Size                 { return nil }
func (m mockOutput) SetResolution(resolution wlc.Size, scale uint32) {}
func (m mockOutput) GetScale() uint32                                { return 0 }
func (m mockOutput) GetSleep() bool                                  { return false }
func (m mockOutput) GetViews() []wlc.View                            { return nil }
func (m mockOutput) ScheduleRender()                                 {}
func (m mockOutput) SetMask(mask uint32)                             {}
func (m mockOutput) SetSleep(sleep bool)                             {}
func (m mockOutput) SetViews(views []wlc.View) bool                  { return false }

// mockOutput mocks the backend.Output interface with a different name than
// mockOutput.
type mockOutput2 struct {
	mockOutput
}

func (m mockOutput2) Name() string {
	return "output2"
}

// TestRootType tests getting the root type from a root container.
func TestRootType(t *testing.T) {
	r := NewRoot()
	if r.Type() != CRoot {
		t.Errorf("expected container type CRoot, got %d", r.Type())
	}
}

// TestRootGeometry tests getting nil geometry from root container.
func TestRootGeometry(t *testing.T) {
	r := NewRoot()
	if r.Geometry() != nil {
		t.Errorf("expected nil value geometry, got %v", r.Geometry())
	}
}

// TestRootChildren tests getting children containers from root container.
func TestRootChildren(t *testing.T) {
	r := NewRoot()
	o := NewOutput(wlc.Output(0), r)
	r.AddChild(o)
	if len(r.Children()) != 1 {
		t.Errorf("expected 1 output, got %d", len(r.Children()))
	}
}

// TestRootFloating tests getting 0 floating child containers from root container.
func TestRootFloating(t *testing.T) {
	r := NewRoot()
	if len(r.Floating()) != 0 {
		t.Errorf("expected 0 floating containers, got %d", len(r.Floating()))
	}
}

// TestRootFocused tests getting focused container from the root container.
func TestRootFocused(t *testing.T) {
	r := NewRoot()
	o := NewOutput(wlc.Output(0), r)
	r.AddChild(o)
	if r.Focused() != o {
		t.Errorf("expected to get output %v, got %v", o, r.Focused())
	}

	r = NewRoot()
	if r.Focused().(*Output) != nil {
		t.Errorf("did not expect to get focused output, got %v", r.Focused())
	}
}

// TestRootParent tests getting no parent container for root container.
func TestRootParent(t *testing.T) {
	r := NewRoot()
	if r.Parent() != nil {
		t.Errorf("did not expect to get parent container, got %v", r.Parent())
	}
}

// TestRootVisible tests if root container is visible.
func TestRootVisible(t *testing.T) {
	r := NewRoot()
	if !r.Visible() {
		t.Errorf("expected root container to be visible")
	}
}

// TestRootAddChild tests adding output to root container.
func TestRootAddChild(t *testing.T) {
	r := NewRoot()
	o := NewOutput(wlc.Output(0), r)
	r.AddChild(o)

	var o2 *Workspace
	r.AddChild(o2)
}

// TestRootSortedWorkspaces tests getting an aggregated list of workspaces
// across all outputs.
func TestRootSortedWorkspaces(t *testing.T) {
	r := NewRoot()

	o1 := NewOutput(mockOutput{}, r)
	r.AddChild(o1)
	w1 := NewWorkspace("2", 2, o1)
	o1.AddChild(w1)

	o2 := NewOutput(mockOutput2{}, r)
	r.AddChild(o2)
	w2 := NewWorkspace("1", 1, o2)
	o2.AddChild(w2)

	ws := r.SortedWorkspaces()

	if len(ws) != 2 {
		t.Errorf("expected 2 workspaces, got %d", len(ws))
	}

	for i, num := range []uint{1, 2} {
		if num != ws[i].Num {
			t.Errorf("expected workspace number %d, got %d", num, ws[i].Num)
		}
	}
}
