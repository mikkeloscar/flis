package i3

import (
	"context"
	"testing"

	"github.com/mikkeloscar/flis/config"
	wlc "github.com/mikkeloscar/go-wlc"
)

// TestNewView tests adding a new view to the layout.
func TestNewView(t *testing.T) {
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
	layout.NewView(ctx, wlc.View(0))
}
