package i3

import (
	"context"
	"testing"

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

	output := wlc.Output(0)
	layout := New()
	layout.NewOutput(ctx, output)
}

// TestFindNextWorkspace tests finding the next available workspace.
func TestFindNextWorkspace(t *testing.T) {
	tests := []struct {
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
	}

	for _, test := range tests {
		name, num := findNextWorkspace(test.current, test.names)
		if name != test.name {
			t.Errorf("expected to get workspace name '%s', got '%s'", test.name, name)
		}

		if num != test.num {
			t.Errorf("expected to get workspace num '%d', got '%d'", test.num, num)
		}
	}
}

// TestFindAvailableWorkspaceNum tests finding available workspace number.
func TestFindAvailableWorkspaceNum(t *testing.T) {
	tests := []struct {
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
	}

	for _, test := range tests {
		num := findAvailableWorkspaceNum(test.ws)
		if num != test.num {
			t.Errorf("expected workspace num %d, got %d", test.num, num)
		}
	}
}
