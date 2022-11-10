package svn

import (
	"github.com/csby/goa/assist"
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
	"sort"
)

func NewUser(log gtype.Log, param *controller.Parameter) *User {
	instance := &User{}
	instance.SetLog(log)
	instance.SetParameter(param)

	return instance
}

type User struct {
	base
}

func (s *User) GetPermissions(ctx gtype.Context, ps gtype.Params) {
	token := s.GetToken(ctx.Token())
	if token == nil {
		ctx.Error(gtype.ErrInternal, "凭证无效")
		return
	}
	if token.Ext == nil {
		ctx.Error(gtype.ErrInternal, "登录用户信息为空")
		return
	}
	user, ok := token.Ext.(*assist.AdEntryUser)
	if !ok || user == nil {
		ctx.Error(gtype.ErrInternal, "登录用户信息无效")
		return
	}

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Svn.Api.Uri.User.Permission
	}

	argument := &model.SvnPermissionID{}
	argument.AccountId = user.SID

	apiData := make([]*model.SvnUserPermission, 0)
	ge := s.callApi(uri, argument, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *User) GetPermissionsDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "用户")
	function := catalog.AddFunction(method, uri, "获取当前登录用户权限列表")
	function.SetOutputDataExample([]*model.SvnUserPermission{
		{
			Repository: "test",
			Path:       "/",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *User) GetAll(ctx gtype.Context, ps gtype.Params) {
	ad := &assist.Ad{}
	if s.Cfg != nil {
		ad.Host = s.Cfg.Ad.Host
		ad.Port = s.Cfg.Ad.Port
		ad.Base = s.Cfg.Ad.Base
		ad.Account = s.Cfg.Ad.Account.Account
		ad.Password = s.Cfg.Ad.Account.Password
	}

	all, err := ad.GetAllUsers()
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	results := make(model.SvnUserArray, 0)
	for k, v := range all {
		result := &model.SvnUser{}
		result.Id = k
		result.Name = v.Name
		result.Account = v.Account

		results = append(results, result)
	}

	ctx.Success(results)
}

func (s *User) GetAllDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "用户")
	function := catalog.AddFunction(method, uri, "获取所有用户列表")
	function.SetOutputDataExample([]*model.SvnUser{
		{
			Id:      gtype.NewGuid(),
			Name:    "管理员",
			Account: "Admin",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *User) GetGroups(ctx gtype.Context, ps gtype.Params) {
	ad := &assist.Ad{}
	parentDN := ""
	if s.Cfg != nil {
		ad.Host = s.Cfg.Ad.Host
		ad.Port = s.Cfg.Ad.Port
		ad.Base = s.Cfg.Ad.Base
		ad.Account = s.Cfg.Ad.Account.Account
		ad.Password = s.Cfg.Ad.Account.Password
		parentDN = s.Cfg.Ad.Root.User
	}

	ous, err := ad.GetOrganizationUnits(parentDN)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	results := make([]*model.SvnUserGroup, 0)
	for i := 0; i < len(ous); i++ {
		ou := ous[i]
		group := &model.SvnUserGroup{
			Name:  ou.Name,
			Users: make(model.SvnUserArray, 0),
		}
		users, ue := ad.GetUsers(ou.DN)
		if ue == nil {
			for j := 0; j < len(users); j++ {
				user := users[j]
				result := &model.SvnUser{}
				result.Id = user.SID
				result.Name = user.Name
				result.Account = user.Account

				group.Users = append(group.Users, result)
			}
			sort.Sort(group.Users)
		}

		results = append(results, group)
	}

	ctx.Success(results)
}

func (s *User) GetGroupsDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "用户")
	function := catalog.AddFunction(method, uri, "获取分组用户列表")
	function.SetOutputDataExample([]*model.SvnUserGroup{
		{
			Name: "分组1",
			Users: []*model.SvnUser{
				{
					Id:      gtype.NewGuid(),
					Name:    "管理员",
					Account: "Admin",
				},
			},
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}
