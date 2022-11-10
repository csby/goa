package auth

import (
	"github.com/csby/goa/controller"
	"github.com/csby/gwsf/gtype"
)

const (
	authCatalogRoot   = "授权服务"
	authCatalogAd     = "域控"
	authCatalogWechat = "微信"
)

type base struct {
	controller.Controller
}

func (s *base) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := s.RootCatalog(doc).AddChild(authCatalogRoot)

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
