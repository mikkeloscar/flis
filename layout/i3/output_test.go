package i3

import (
	"context"
	"fmt"
	"testing"

	"github.com/mikkeloscar/flis/backend"
	"github.com/mikkeloscar/flis/config"
	"github.com/mikkeloscar/flis/layout"
	wlc "github.com/mikkeloscar/go-wlc"
)

// TestNewOutput tests adding a new output to the layout.
func TestNewOutput(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		"config",
		&config.Config{
			Workspaces: map[uint]string{
				1: "one",
			},
		},
	)

	output := mockOutput(0)
	layout := New()
	layout.NewOutput(ctx, output)
}

type mockOutput int

func (o mockOutput) Focus()                                          {}
func (o mockOutput) GetMask() uint32                                 { return 0 }
func (o mockOutput) GetMutableViews() []wlc.View                     { return nil }
func (o mockOutput) Name() string                                    { return fmt.Sprintf("%d", o) }
func (o mockOutput) GetRenderer() wlc.Renderer                       { return wlc.NoRenderer }
func (o mockOutput) GetResolution() *wlc.Size                        { return nil }
func (o mockOutput) GetVirtualResolution() *wlc.Size                 { return nil }
func (o mockOutput) SetResolution(resolution wlc.Size, scale uint32) {}
func (o mockOutput) GetScale() uint32                                { return 0 }
func (o mockOutput) GetSleep() bool                                  { return false }
func (o mockOutput) GetViews() []wlc.View                            { return nil }
func (o mockOutput) ScheduleRender()                                 {}
func (o mockOutput) SetMask(mask uint32)                             {}
func (o mockOutput) SetSleep(sleep bool)                             {}
func (o mockOutput) SetViews(views []wlc.View) bool                  { return false }

// TestOutputByBackend tests getting an output container from a backend output
// interface.
func TestOutputByBackend(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		"config",
		&config.Config{
			Workspaces: map[uint]string{
				1: "one",
			},
		},
	)

	layout := New()

	for i := 0; i < 5; i++ {
		output := mockOutput(i)
		layout.NewOutput(ctx, output)
	}

	for _, ti := range []struct {
		output backend.Output
		exists bool
	}{
		{
			output: mockOutput(2),
			exists: true,
		},
		{
			output: mockOutput(10),
			exists: false,
		},
	} {
		c := layout.OutputByBackend(ti.output)
		if c == nil && ti.exists {
			t.Errorf("expected to find an output")
		}

		if c != nil && !ti.exists {
			t.Errorf("did not expect to find an output")
		}
	}
}

// TestFindNextWorkspace tests finding the next available workspace.
func TestFindNextWorkspace(t *testing.T) {
	for _, ti := range []struct {
		current []*layout.Workspace
		names   map[uint]string
		name    string
		num     uint
	}{
		{
			[]*layout.Workspace{
				{Num: 1},
				{Num: 2},
			},
			map[uint]string{
				2: "two",
			},
			"3",
			3,
		},
		{
			[]*layout.Workspace{
				{Num: 1},
				{Num: 2},
			},
			map[uint]string{
				3: "three",
			},
			"three",
			3,
		},
		{
			[]*layout.Workspace{
				{Num: 1},
			},
			nil,
			"2",
			2,
		},
	} {
		name, num := findNextWorkspace(ti.current, ti.names)
		if name != ti.name {
			t.Errorf("expected to get workspace name '%s', got '%s'", ti.name, name)
		}

		if num != ti.num {
			t.Errorf("expected to get workspace num '%d', got '%d'", ti.num, num)
		}
	}
}

// TestFindAvailableWorkspaceNum tests finding available workspace number.
func TestFindAvailableWorkspaceNum(t *testing.T) {
	for _, ti := range []struct {
		ws  []*layout.Workspace
		num uint
	}{
		{
			[]*layout.Workspace{
				{Num: 2},
				{Num: 5},
				{Num: 7},
			},
			1,
		},
		{

			[]*layout.Workspace{
				{Num: 1},
				{Num: 2},
			},
			3,
		},
		{

			[]*layout.Workspace{
				{Num: 1},
				{Num: 3},
			},
			2,
		},
	} {
		num := findAvailableWorkspaceNum(ti.ws)
		if num != ti.num {
			t.Errorf("expected workspace num %d, got %d", ti.num, num)
		}
	}
}
