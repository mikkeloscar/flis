package i3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mikkeloscar/flise/layout"
)

type Layout struct {
	root *layout.Root
}

func New() *Layout {
	return &Layout{
		root: layout.NewRoot(),
	}
}

func (l *Layout) Arrange(start layout.Container) {
	// TODO: implement with more than one view
	switch c := start.(type) {
	case *layout.View:
		pG := c.Parent().Geometry()
		g := c.Geometry()
		g.Origin = pG.Origin
		g.Size = pG.Size
		log.Debugf("Arranging view %s %dx%d (%d,%d)", c.Title(),
			g.Size.W, g.Size.H, g.Origin.X, g.Origin.Y)
	}
}

func (l *Layout) ArrangeRoot() {
	l.Arrange(l.root)
}
