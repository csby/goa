package ad

import (
	"github.com/csby/goa/controller"
	"github.com/csby/gwsf/gtype"
)

const (
	adCatalogRoot   = "域控"
	adCatalogUser   = "用户"
	adCatalogGroup  = "组"
	adCatalogServer = "服务器"
	adCatalogShare  = "共享目录"
)

type base struct {
	controller.Controller
}

func (s *base) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := s.RootCatalog(doc).AddChild(adCatalogRoot)

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
