package dhcp

import (
	"github.com/csby/goa/controller"
	"github.com/csby/gwsf/gtype"
)

type base struct {
	controller.Controller
}

func (s *base) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := s.RootCatalog(doc).AddChild("DHCP")

	count := len(names)
	if count < 1 {
		return root
	}

	child := root
	for i := 0; i < count; i++ {
		name := names[i]
		child = child.AddChild(name)
	}

	return child
}
