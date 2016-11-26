package layout

import (
	"context"
	"testing"
)

// TestGetLayout tests getting layout from context.
func TestGetLayout(t *testing.T) {
	ctx := context.WithValue(context.Background(), "layout", Mock{})

	layout := Get(ctx)
	if layout == nil {
		t.Errorf("expected to get layout, got nil")
	}
}
