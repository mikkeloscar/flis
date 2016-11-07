package layout

import (
	"github.com/mikkeloscar/flis/backend"
	"github.com/mikkeloscar/flis/context"
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
	NewWorkspace(ctx context.Context, output *Output, name string)

	// NewView initializes a new view.
	NewView(ctx context.Context, view backend.View)

	// ArrangeRoot arranges the whole layout from the root and down.
	ArrangeRoot()

	// Arrange arranges a subbranch of the layout starting from the
	// specified container and moving down the layout.
	Arrange(start Container)
}

// Get layout from context.
func Get(ctx context.Context) *Layout {
	return ctx.MustGet("layout").(*Layout)
}

type ContainerType int

const (
	CRoot ContainerType = iota
	COutput
	CWorkspace
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
	Fullscreen() Container

	// Parent of the container.
	Parent() Container

	// Visible returns true if the container is visible.
	Visible() bool

	AddChild(Container)

	Type() ContainerType
}
