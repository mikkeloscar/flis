package layout

import wlc "github.com/mikkeloscar/go-wlc"

type Workspace struct {
	name       string
	containers []Container
	floating   []Container
	focused    Container
	output     *Output
}

func NewWorkspace(name string, output *Output) *Workspace {
	return &Workspace{
		name:       name,
		containers: make([]Container, 0),
		floating:   make([]Container, 0),
		focused:    nil,
		output:     output,
	}
}

func (w *Workspace) Type() ContainerType {
	return CWorkspace
}

func (w *Workspace) Name() string {
	return w.name
}

// Geometry for the workspace is the geometry of the parent output.
func (w *Workspace) Geometry() *wlc.Geometry {
	return w.output.Geometry()
}

func (w *Workspace) Children() []Container {
	return w.containers
}

func (w *Workspace) Floating() []Container {
	return w.floating
}

func (w *Workspace) Focused() Container {
	return w.focused
}

func (w *Workspace) Fullscreen() Container {
	if w.focused != nil {
		return w.focused.Fullscreen()
	}
	return nil
}

func (w *Workspace) Parent() Container {
	return w.output
}

func (w *Workspace) AddChild(container Container) {
	// TODO:
}

func (w *Workspace) Visible() bool {
	return w.output.Visible() && w.output.focused == w
}
