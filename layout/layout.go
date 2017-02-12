package layout

import (
	"context"

	"github.com/mikkeloscar/flis/backend"
	wlc "github.com/mikkeloscar/go-wlc"
)

// Direction defines a movement direction in a layout.
type Direction int

const (
	// Up is direction up.
	Up Direction = iota
	// Down is direction down.
	Down
	// Left is direction left.
	Left
	// Right is direction right.
	Right
)

// Layout defines an interface for interacting with a layout.
type Layout interface {
	// Move container in direction.
	Move(ctx context.Context, dir Direction)

	// Focus container.
	Focus(ctx context.Context, container Container)

	// Focused returns the focused container of the layout.
	Focused(ctx context.Context) Container

	FocusedByType(ctx context.Context, typ ContainerType) Container

	// NewOuput initializes a new output.
	NewOutput(ctx context.Context, output backend.Output)

	// NewWorkspace initializes a new workspace.
	NewWorkspace(ctx context.Context, output *Output, name string, num uint)

	// NewView initializes a new view.
	NewView(ctx context.Context, view backend.View) Container

	// ArrangeRoot arranges the whole layout from the root and down.
	ArrangeRoot()

	// Arrange arranges a subbranch of the layout starting from the
	// specified container and moving down the layout.
	Arrange(start Container)

	// OutputByBackend gets output container from backend output interface.
	OutputByBackend(output backend.Output) Container
	// ViewByBackend(output backend.Output) Container
}

// Get layout from context.
func Get(ctx context.Context) Layout {
	return ctx.Value("layout").(Layout)
}

// ContainerType is a type of container e.g. Root or Output.
type ContainerType int

const (
	// CRoot is the root container type.
	CRoot ContainerType = iota
	// COutput is the output container type.
	COutput
	// CWorkspace is the workspace container type.
	CWorkspace
	// CView is the view container type.
	CView
)

// Container defines an interface describing a container e.g. output,
// workspace or view.
type Container interface {
	// Geometry of the container.
	Geometry() *wlc.Geometry

	// Children of the container.
	Children() []Container

	// Floating children of the container.
	Floating() []Container

	// Focused child of the container.
	Focused() Container

	// // IsFocused is true if the container has focus.
	// IsFocused() bool
	// SetFocused set focused child of the container.
	// SetFocused(Container)
	// Fullscreen returns the fullscreened view of a container, if any.
	// Fullscreen() Container

	// Parent of the container.
	Parent() Container

	// Visible returns true if the container is visible.
	Visible() bool

	// AddChild adds a child container.
	AddChild(Container)

	// Type returns the type of a container.
	Type() ContainerType
}
