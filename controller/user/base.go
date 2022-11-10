package user

import (
	"github.com/csby/goa/controller"
	"github.com/csby/gwsf/gtype"
)

const (
	userCatalogRoot  = "用户管理"
	userCatalogLogin = "登录用户"
)

type base struct {
	controller.Controller
}

func (s *base) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := s.RootCatalog(doc).AddChild(userCatalogRoot)

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
