package svn

import (
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
)

func NewPermission(log gtype.Log, param *controller.Parameter) *Permission {
	instance := &Permission{}
	instance.SetLog(log)
	instance.SetParameter(param)

	return instance
}

type Permission struct {
	base
}

func (s *Permission) GetItemList(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnRepository{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(name)为空")
		return
	}
	if len(argument.Path) < 1 {
		ctx.Error(gtype.ErrInput, "路径(path)为空")
		return
	}

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Svn.Api.Uri.Permission.List
	}

	apiData := make([]*model.SvnPermission, 0)
	ge := s.callApi(uri, argument, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Permission) GetItemListDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "权限")
	function := catalog.AddFunction(method, uri, "获取项目访问权限列表")
	function.SetInputJsonExample(&model.SvnRepository{
		Name: "dd",
		Path: "/trunk",
	})
	function.SetOutputDataExample([]*model.SvnPermission{
		{},
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Permission) AddItem(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnPermissionArgumentEdit{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Repository) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(repository)为空")
		return
	}
	if len(argument.Path) < 1 {
		ctx.Error(gtype.ErrInput, "路径(path)为空")
		return
	}
	if len(argument.AccountId) < 1 {
		ctx.Error(gtype.ErrInput, "账号ID(accountId)为空")
		return
	}

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Svn.Api.Uri.Permission.Add
	}

	ge := s.callApi(uri, argument, nil)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(nil)
}

func (s *Permission) AddItemDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "权限")
	function := catalog.AddFunction(method, uri, "添加项目访问权限")
	function.SetInputJsonExample(&model.SvnPermissionArgumentEdit{
		SvnPermissionArgument: model.SvnPermissionArgument{
			Repository: "test",
			Path:       "/trunk",
			AccountId:  "S-1-5-21-1114322273-403004966-1807125474-500",
		},
		AccessLevel: 2,
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Permission) ModItem(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnPermissionArgumentEdit{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Repository) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(repository)为空")
		return
	}
	if len(argument.Path) < 1 {
		ctx.Error(gtype.ErrInput, "路径(path)为空")
		return
	}
	if len(argument.AccountId) < 1 {
		ctx.Error(gtype.ErrInput, "账号ID(accountId)为空")
		return
	}

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Svn.Api.Uri.Permission.Mod
	}

	ge := s.callApi(uri, argument, nil)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(nil)
}

func (s *Permission) ModItemDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "权限")
	function := catalog.AddFunction(method, uri, "修改项目访问权限")
	function.SetInputJsonExample(&model.SvnPermissionArgumentEdit{
		SvnPermissionArgument: model.SvnPermissionArgument{
			Repository: "test",
			Path:       "/trunk",
			AccountId:  "S-1-5-21-1114322273-403004966-1807125474-500",
		},
		AccessLevel: 2,
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Permission) DelItem(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnPermissionArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Repository) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(repository)为空")
		return
	}
	if len(argument.Path) < 1 {
		ctx.Error(gtype.ErrInput, "路径(path)为空")
		return
	}
	if len(argument.AccountId) < 1 {
		ctx.Error(gtype.ErrInput, "账号ID(accountId)为空")
		return
	}

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Svn.Api.Uri.Permission.Del
	}

	ge := s.callApi(uri, argument, nil)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(nil)
}

func (s *Permission) DelItemDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "权限")
	function := catalog.AddFunction(method, uri, "删除项目访问权限")
	function.SetInputJsonExample(&model.SvnPermissionArgument{
		Repository: "test",
		Path:       "/trunk",
		AccountId:  "S-1-5-21-1114322273-403004966-1807125474-500",
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}
