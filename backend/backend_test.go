package backend

import (
	"context"
	"testing"
)

func TestMustGetBackend(t *testing.T) {
	ctx := context.WithValue(context.Background(), "backend", Mock{})
	Get(ctx)
}
