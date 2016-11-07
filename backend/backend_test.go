package backend

import (
	"testing"

	"github.com/mikkeloscar/flis/context"
)

func TestMustGetBackend(t *testing.T) {
	ctx := context.Context(map[string]interface{}{
		"backend": Mock{},
	})

	Get(ctx)
}
