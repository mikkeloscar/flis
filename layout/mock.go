package layout

import (
	"context"

	"github.com/mikkeloscar/flis/backend"
	wlc "github.com/mikkeloscar/go-wlc"
)

// Mock mocks the layout interface.
type Mock struct{}

// Move mocks moving container in direction.
func (m Mock) Move(ctx context.Context, dir Direction) {}

// Focus mocks focusing a container.
func (m Mock) Focus(ctx context.Context, container Container) {}

// Focused mocks returning the focused container of the layout.
func (m Mock) Focused(ctx context.Context) Container { return nil }

// FocusedByType mocks returning a focused container by type.
func (m Mock) FocusedByType(ctx context.Context, typ ContainerType) Container { return nil }

// NewOutput mocks initializing a new output.
func (m Mock) NewOutput(ctx context.Context, output backend.Output) {}

// NewWorkspace mocks initializing a new workspace.
func (m Mock) NewWorkspace(ctx context.Context, output *Output, name string, num uint) {}

// NewView mocks initializing a new view.
func (m Mock) NewView(ctx context.Context, view backend.View) Container {
	return MockView{}
}

// ArrangeRoot mocks arranging the whole layout from the root and down.
func (m Mock) ArrangeRoot() {}

// Arrange mocks arranging a subbranch of the layout starting from the
// specified container and moving down the layout.
func (m Mock) Arrange(start Container) {}

// OutputByBackend gets output container from backend output interface.
func (m Mock) OutputByBackend(output backend.Output) Container { return nil }

// MockView mocks a view container.
type MockView struct{}

// Geometry mocks geometry of the container.
func (v MockView) Geometry() *wlc.Geometry { return nil }

// Children mocks children of the container.
func (v MockView) Children() []Container { return nil }

// Floating mocks floating children of the container.
func (v MockView) Floating() []Container { return nil }

// Focused mocks focused child of the container.
func (v MockView) Focused() Container { return nil }

// Parent mocks getting parent container of the container.
func (v MockView) Parent() Container { return nil }

// Visible mocks return the visibility state of the container.
func (v MockView) Visible() bool { return false }

// AddChild mocks adding a child container.
func (v MockView) AddChild(Container) {}

// Type mocks returning the type of the container.
func (v MockView) Type() ContainerType { return CView }
