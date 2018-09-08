package i3

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/mikkeloscar/flis/layout"
)

// Focused returns the focused container lowest in the three, or nil if none is
// found.
func (l *Layout) Focused(ctx context.Context) layout.Container {
	var current layout.Container
	current = l.root
	for current != nil {
		if current.Focused() == nil {
			return current
		}

		current = current.Focused()
	}

	return nil
}

// Focus sets focus to a container.
func (l *Layout) Focus(ctx context.Context, c layout.Container) {
	// TODO: set focus to a container
	// return true

}

// FocusedByType returns the currently focused container by type.
func (l *Layout) FocusedByType(ctx context.Context, typ layout.ContainerType) layout.Container {
	var current layout.Container
	current = l.root
	for current != nil {
		if current.Type() == typ {
			return current
		}

		current = current.Focused()
	}

	// TODO: translate type to string
	log.Warnf("Failed to find focused container of type %d", typ)
	return nil
}
