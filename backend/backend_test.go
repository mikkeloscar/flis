package backend

import (
	"testing"

	"github.com/mikkeloscar/flise/context"
)

func TestMustGetBackend(t *testing.T) {
	ctx := context.Context(map[string]interface{}{
		"backend": Mock{},
	})

	Get(ctx)
}
