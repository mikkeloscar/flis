package layout

import (
	"github.com/mikkeloscar/flis/backend"
	"github.com/mikkeloscar/flis/context"
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
func (m Mock) NewWorkspace(ctx context.Context, output *Output, name string) {}

// NewView mocks initializing a new view.
func (m Mock) NewView(ctx context.Context, view backend.View) {}

// ArrangeRoot mocks arranging the whole layout from the root and down.
func (m Mock) ArrangeRoot() {}

// Arrange mocks arranging a subbranch of the layout starting from the
// specified container and moving down the layout.
func (m Mock) Arrange(start Container) {}
