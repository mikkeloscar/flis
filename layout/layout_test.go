package layout

import (
	"testing"

	"github.com/mikkeloscar/flis/context"
)

// TestGetLayout tests getting layout from context.
func TestGetLayout(t *testing.T) {
	ctx := context.Context(map[string]interface{}{
		"layout": Mock{},
	})

	layout := Get(ctx)
	if layout == nil {
		t.Errorf("expected to get layout, got nil")
	}
}
