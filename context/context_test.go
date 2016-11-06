package context

import "testing"

// TestSetContext tests setting a value in the context.
func TestSetContext(t *testing.T) {
	ctx := Context(map[string]interface{}{})
	ctx.Set("key", "value")
}

// TestMustGetContext tests getting value from context.
func TestMustGetContext(t *testing.T) {
	ctx := Context(map[string]interface{}{
		"key": "value",
	})

	v := ctx.MustGet("key")

	if s, ok := v.(string); !ok || s != "value" {
		t.Errorf("should get correct value")
	}

	// recover from panic in MustGet call
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected to recover from panic")
		}
	}()

	ctx.MustGet("invalid_key")
}
