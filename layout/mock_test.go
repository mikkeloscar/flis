package layout

import "testing"

// TestMock creates test coverage for the layout mock.
func TestMock(t *testing.T) {
	m := Mock{}

	m.Move(nil, Up)
	m.Focus(nil, nil)
	m.Focused(nil)
	m.FocusedByType(nil, CRoot)
	m.NewOutput(nil, nil)
	m.NewWorkspace(nil, nil, "", 0)
	m.NewView(nil, nil)
	m.ArrangeRoot()
	m.Arrange(nil)
	m.OutputByBackend(nil)

	v := MockView{}

	v.Geometry()
	v.Children()
	v.Floating()
	v.Focused()
	v.Parent()
	v.Visible()
	v.AddChild(nil)
	v.Type()
}
