package layout

import (
	"testing"

	wlc "github.com/mikkeloscar/go-wlc"
)

// TestWorkspaceType tests getting the workspace type from a workspace container.
func TestWorkspaceType(t *testing.T) {
	w := NewWorkspace("1", 1, nil)
	if w.Type() != CWorkspace {
		t.Errorf("expected container type CWorkspace, got %d", w.Type())
	}
}

// TestWorkspaceGeometry tests getting geometry of workspace which should be
// the same as the output.
func TestWorkspaceGeometry(t *testing.T) {
	o := NewOutput(mockOutputRes{}, nil)
	w := NewWorkspace("1", 1, o)
	o.AddChild(w)
	if !wlc.GeometryEquals(*w.Geometry(), wlc.GeometryZero) {
		t.Errorf("expected geometry (0,0 - 0x0), got %v", w.Geometry())
	}
}

// TestWorkspaceChildren tests getting children containers of the workspace container.
func TestWorkspaceChildren(t *testing.T) {
	w := NewWorkspace("1", 1, nil)
	v := NewView(wlc.View(0), w)
	w.AddChild(v)
	if len(w.Children()) != 1 {
		t.Errorf("expected 1 child container, got %d", len(w.Children()))
	}
}

// TestWorkspaceFloating tests getting 0 floating child containers from workspace container.
func TestWorkspaceFloating(t *testing.T) {
	w := NewWorkspace("1", 1, nil)
	if len(w.Floating()) != 0 {
		t.Errorf("expected 0 floating containers, got %d", len(w.Floating()))
	}
}

// TestWorkspaceFocused tests getting focused container from the workspace container.
func TestWorkspaceFocused(t *testing.T) {
	w := NewWorkspace("1", 1, nil)
	v := NewView(wlc.View(0), w)
	w.AddChild(v)
	if w.Focused() != v {
		t.Errorf("expected to get container %v, got %v", v, w.Focused())
	}

	w = NewWorkspace("1", 1, nil)
	if w.Focused() != nil {
		t.Errorf("did not expect to get focused container, got %v", w.Focused())
	}
}

// TestWorkspaceParent tests getting parent output for workspace container.
func TestWorkspaceParent(t *testing.T) {
	o := NewOutput(mockOutputRes{}, nil)
	w := NewWorkspace("1", 1, o)
	o.AddChild(w)
	if w.Parent() != o {
		t.Errorf("expected to get parent %v, got %v", o, w.Parent())
	}
}

// TestWorkspaceVisible tests if workspace is visible.
func TestWorkspaceVisible(t *testing.T) {
	o := NewOutput(mockOutputRes{}, nil)
	w := NewWorkspace("1", 1, o)
	o.AddChild(w)
	if !w.Visible() {
		t.Errorf("expected workspace container to be visible")
	}

	w2 := NewWorkspace("2", 2, o)
	o.AddChild(w2)
	if w2.Visible() {
		t.Errorf("did not expect workspace container to be visible")
	}
}

// TestWorkspaceAddChild tests adding view to a workspace.
func TestWorkspaceAddChild(t *testing.T) {
	w := NewWorkspace("1", 1, nil)
	v := NewView(wlc.View(0), w)
	w.AddChild(v)

	var v2 *Workspace
	w.AddChild(v2)
}

// TestWorkspacesLen tests getting correct length of workspace list.
func TestWorkspacesLen(t *testing.T) {
	workspaces := []struct {
		workspaces Workspaces
		len        int
	}{
		{
			Workspaces([]*Workspace{}),
			0,
		},
		{
			Workspaces([]*Workspace{
				{},
			}),
			1,
		},
	}

	for _, w := range workspaces {
		if w.len != w.workspaces.Len() {
			t.Errorf("expected %d, got %d", w.len, w.workspaces.Len())
		}
	}
}

// TestWorkspacesSwap tests swapping workspaces in workspace list.
func TestWorkspacesSwap(t *testing.T) {
	workspaces := Workspaces([]*Workspace{
		{Num: 1},
		{Num: 2},
	})

	workspaces.Swap(0, 1)

	if workspaces[0].Num != 2 || workspaces[1].Num != 1 {
		t.Errorf("expected workspaces to be swapped")
	}
}

// TestWorkspacesLess tests workspaces less function.
func TestWorkspacesLess(t *testing.T) {
	workspaces := Workspaces([]*Workspace{
		{Num: 1},
		{Num: 2},
	})

	if !workspaces.Less(0, 1) {
		t.Errorf("expected workspace 0 to be less than workspace 1")
	}
}
