package layout

import (
	"testing"

	wlc "github.com/mikkeloscar/go-wlc"
)

// mockOutputRes mocks the backend.Output interface with a size of 0x0.
type mockOutputRes struct {
	mockOutput
}

func (m mockOutputRes) GetResolution() *wlc.Size {
	return &wlc.SizeZero
}

func (m mockOutputRes) GetVirtualResolution() *wlc.Size {
	return &wlc.SizeZero
}

// TestOutputType tests getting the output type from an output.
func TestOutputType(t *testing.T) {
	o := NewOutput(wlc.Output(0), nil)
	if o.Type() != COutput {
		t.Errorf("expected container type COutput, got %d", o.Type())
	}
}

// TestOutputGeometry tests getting geometry of output.
func TestOutputGeometry(t *testing.T) {
	o := NewOutput(mockOutputRes{}, nil)
	if !wlc.GeometryEquals(*o.Geometry(), wlc.GeometryZero) {
		t.Errorf("expected geometry (0,0 - 0x0), got %v", o.Geometry())
	}
}

// TestOutputChildren tests getting children containers of an output.
func TestOutputChildren(t *testing.T) {
	o := NewOutput(wlc.Output(0), nil)
	w := NewWorkspace("1", 1, o)
	o.AddChild(w)
	if len(o.Children()) != 1 {
		t.Errorf("expected 1 child container, got %d", len(o.Children()))
	}
}

// TestOutputFloating tests getting 0 floating child containers from an output.
func TestOutputFloating(t *testing.T) {
	o := NewOutput(wlc.Output(0), nil)
	if len(o.Floating()) != 0 {
		t.Errorf("expected 0 floating containers, got %d", len(o.Floating()))
	}
}

// TestOutputFocused tests getting focused container from an output.
func TestOutputFocused(t *testing.T) {
	o := NewOutput(wlc.Output(0), nil)
	w := NewWorkspace("1", 1, o)
	o.AddChild(w)
	if o.Focused() != w {
		t.Errorf("expected to get container %v, got %v", w, o.Focused())
	}

	o = NewOutput(wlc.Output(0), nil)
	if o.Focused().(*Workspace) != nil {
		t.Errorf("did not expect to get focused container, got %v", o.Focused())
	}
}

// TestOutputParent tests getting root container an output.
func TestOutputParent(t *testing.T) {
	r := NewRoot()
	o := NewOutput(mockOutputRes{}, r)
	r.AddChild(o)
	if o.Parent() != r {
		t.Errorf("expected to get parent %v, got %v", r, o.Parent())
	}
}

// TestOutputVisible tests if output is visible.
func TestOutputVisible(t *testing.T) {
	o := NewOutput(mockOutputRes{}, nil)
	if !o.Visible() {
		t.Errorf("expected output container to be visible")
	}
}

// TestOutputAddChild tests adding workspace to an output.
func TestOutputAddChild(t *testing.T) {
	o := NewOutput(mockOutputRes{}, nil)
	w := NewWorkspace("1", 1, o)
	o.AddChild(w)

	var w2 *Output
	o.AddChild(w2)
}
