package ad

import (
	"fmt"
	"github.com/csby/goa/assist"
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
	"sort"
	"strings"
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

func (s *User) CreateUser(ctx gtype.Context, ps gtype.Params) {
	token := s.GetToken(ctx.Token())
	if token == nil {
		ctx.Error(gtype.ErrInternal, "凭证无效")
		return
	}
	if !s.IsAdmin(token.UserAccount) {
		ctx.Error(gtype.ErrNoPermission, "需要管理员权限才能新建用户")
		return
	}

	argument := &model.AdUserCreate{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	user := &assist.AdEntryUserCreate{}
	user.Name = strings.TrimSpace(argument.Name)
	if len(user.Name) < 1 {
		ctx.Error(gtype.ErrInput, "用户姓名(name)为空")
		return
	}
	user.Account = strings.TrimSpace(argument.Account)
	if len(user.Account) < 1 {
		ctx.Error(gtype.ErrInput, "登录帐号(account)为空")
		return
	}
	user.Password = strings.TrimSpace(argument.Password)
	if len(user.Password) < 1 {
		ctx.Error(gtype.ErrInput, "登录密码(password)为空")
		return
	}
	if len(argument.Manager) > 0 {
		user.Manager, err = s.FromBase64(argument.Manager)
		if err != nil {
			ctx.Error(gtype.ErrInput, fmt.Errorf("直接主管DN(manager)无效: %s", err.Error()))
			return
		}
	}
	if len(argument.Parent) > 0 {
		user.Parent, err = s.FromBase64(argument.Parent)
		if err != nil {
			ctx.Error(gtype.ErrInput, fmt.Errorf("组织单位DN(parent)无效: %s", err.Error()))
			return
		}
	}

	ad := s.Ad()
	result, err := ad.NewUser(user)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	ctx.Success(&model.AdUser{
		AdDn: model.AdDn{
			Dn: s.ToBase64(result.DN),
		},
		SID:     result.SID,
		Name:    result.Name,
		Account: result.Account,
	})
}

func (s *User) CreateUserDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogUser)
	function := catalog.AddFunction(method, uri, "新建用户")
	function.SetNote("成功时返回新建用户的信息")
	function.SetInputJsonExample(&model.AdUserCreate{
		Account: "admin",
		Name:    "管理员",
	})
	function.SetOutputDataExample(&model.AdUser{
		Account: "admin",
		Name:    "管理员",
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
}

func (s *User) ResetPassword(ctx gtype.Context, ps gtype.Params) {
	argument := &model.AdSetPassword{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	account := strings.TrimSpace(argument.Account)
	if len(account) < 1 {
		ctx.Error(gtype.ErrInput, "登录帐号(account)为空")
		return
	}

	ad := s.Ad()
	err = ad.SetUserPassword(account, argument.Password)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	ctx.Success(nil)
}

func (s *User) ResetPasswordDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogUser)
	function := catalog.AddFunction(method, uri, "重置密码")
	function.SetInputJsonExample(&model.AdSetPassword{})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInput)
}

func (s *User) ChangePassword(ctx gtype.Context, ps gtype.Params) {
	token := s.GetToken(ctx.Token())
	if token == nil {
		ctx.Error(gtype.ErrInternal, "凭证无效")
		return
	}

	argument := &model.AdChangePassword{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Account) < 1 {
		argument.Account = token.UserAccount
	}

	account := strings.TrimSpace(argument.Account)
	if len(account) < 1 {
		ctx.Error(gtype.ErrInput, "登录帐号(account)为空")
		return
	}

	ad := s.Ad()
	err = ad.ChangeUserPassword(account, argument.OldPassword, argument.NewPassword)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	ctx.Success(nil)
}

func (s *User) ChangePasswordDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogUser)
	function := catalog.AddFunction(method, uri, "修改密码")
	function.SetNote("未指定帐号时，修改当前登录用户的密码")
	function.SetInputJsonExample(&model.AdChangePassword{})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInput)
}

func (s *User) GetAccountList(ctx gtype.Context, ps gtype.Params) {
	ad := s.Ad()
	users, err := ad.GetAllUsers()
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	results := make(model.AdUserCollection, 0)
	for k, v := range users {
		results = append(results, &model.AdUser{
			AdDn: model.AdDn{
				Dn: s.ToBase64(v.DN),
			},
			SID:     k,
			Account: v.Account,
			Name:    v.Name,
		})
	}

	sort.Sort(results)
	ctx.Success(results)
}

func (s *User) GetAccountListDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogUser)
	function := catalog.AddFunction(method, uri, "获取帐号列表")
	function.SetNote("获取所有用户帐号信息列表")
	function.SetOutputDataExample(model.AdUserCollection{
		{
			Name: "张三",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *User) GetAccountTree(ctx gtype.Context, ps gtype.Params) {
	ad := s.Ad()
	root := s.Cfg.Ad.Root.User
	units, err := ad.GetOrganizationUnits(root)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	results := make(model.AdOrganizationUnitCollection, 0)
	c := len(units)
	for i := 0; i < c; i++ {
		item := units[i]
		if item == nil {
			continue
		}

		result := &model.AdOrganizationUnit{}
		result.Dn = s.ToBase64(item.DN)
		result.Name = item.Name
		result.Description = item.Description
		result.Street = item.Street
		result.Users = make(model.AdUserCollection, 0)

		users, ue := ad.GetUsers(item.DN)
		if ue == nil {
			uc := len(users)
			for ui := 0; ui < uc; ui++ {
				user := users[ui]
				if user == nil {
					continue
				}

				result.Users = append(result.Users, &model.AdUser{
					AdDn: model.AdDn{
						Dn: s.ToBase64(user.DN),
					},
					SID:     user.SID,
					Account: user.Account,
					Name:    user.Name,
				})
			}
			sort.Sort(result.Users)
		}

		results = append(results, result)
	}

	sort.Sort(results)
	ctx.Success(results)
}

func (s *User) GetAccountTreeDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogUser)
	function := catalog.AddFunction(method, uri, "获取帐号树")
	function.SetNote("获取帐号及分组信息信息")
	function.SetOutputDataExample([]*model.AdOrganizationUnit{
		{
			Name: "开发部",
			Users: model.AdUserCollection{
				{
					Name: "张三",
				},
			},
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *User) GetOrganizationUnitList(ctx gtype.Context, ps gtype.Params) {
	ad := s.Ad()
	root := s.Cfg.Ad.Root.User
	units, err := ad.GetOrganizationUnits(root)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	results := make(model.AdOrganizationUnitCollection, 0)
	c := len(units)
	for i := 0; i < c; i++ {
		item := units[i]
		if item == nil {
			continue
		}

		result := &model.AdOrganizationUnit{}
		result.Dn = s.ToBase64(item.DN)
		result.Name = item.Name
		result.Description = item.Description
		result.Street = item.Street

		results = append(results, result)
	}

	sort.Sort(results)
	ctx.Success(results)
}

func (s *User) GetOrganizationUnitListDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogUser)
	function := catalog.AddFunction(method, uri, "获取组织单位列表")
	function.SetOutputDataExample([]*model.AdOrganizationUnit{
		{
			Name: "开发部",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *User) GetSubordinates(ctx gtype.Context, ps gtype.Params) {
	token := s.GetToken(ctx.Token())
	if token == nil {
		ctx.Error(gtype.ErrInternal, "凭证无效")
		return
	}

	argument := &model.AdAccount{}
	ctx.GetJson(argument)
	if len(argument.Account) < 1 {
		argument.Account = token.UserAccount
	}

	ad := s.Ad()
	users, err := ad.GetUserSubordinates(argument.Account)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	results := make(model.AdUserCollection, 0)
	c := len(users)
	for i := 0; i < c; i++ {
		u := users[i]
		if u == nil {
			continue
		}
		results = append(results, &model.AdUser{
			AdDn: model.AdDn{
				Dn: s.ToBase64(u.DN),
			},
			SID:     u.SID,
			Account: u.Account,
			Name:    u.Name,
		})
	}

	sort.Sort(results)
	ctx.Success(results)
}

func (s *User) GetSubordinatesDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogUser)
	function := catalog.AddFunction(method, uri, "获取下属列表")
	function.SetNote("如果未指定帐号，默认为当前登录用户")
	function.SetInputJsonExample(&model.AdAccount{})
	function.SetOutputDataExample(&model.AdUserCollection{
		{
			Name: "张三",
		},
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
}

func (s *User) GetVpnEnable(ctx gtype.Context, ps gtype.Params) {
	token := s.GetToken(ctx.Token())
	if token == nil {
		ctx.Error(gtype.ErrInternal, "凭证无效")
		return
	}

	argument := &model.AdAccount{}
	ctx.GetJson(argument)
	if len(argument.Account) < 1 {
		argument.Account = token.UserAccount
	}

	ad := s.Ad()
	enable, err := ad.GetUserVpnEnable(argument.Account)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	ctx.Success(enable)
}

func (s *User) GetVpnEnableDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogUser)
	function := catalog.AddFunction(method, uri, "获取VPN启用状态")
	function.SetNote("如果未指定帐号，默认为当前登录用户")
	function.SetInputJsonExample(&model.AdAccount{})
	function.SetOutputDataExample(true)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
}

func (s *User) SetVpnEnable(ctx gtype.Context, ps gtype.Params) {
	token := s.GetToken(ctx.Token())
	if token == nil {
		ctx.Error(gtype.ErrInternal, "凭证无效")
		return
	}

	argument := &model.AdVpn{}
	ctx.GetJson(argument)
	if len(argument.Account) < 1 {
		argument.Account = token.UserAccount
	}

	ad := s.Ad()
	err := ad.SetUserVpnEnable(argument.Account, argument.Enable)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	ctx.Success(nil)
}

func (s *User) SetVpnEnableDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogUser)
	function := catalog.AddFunction(method, uri, "设置VPN启用状态")
	function.SetNote("如果未指定帐号，默认为当前登录用户")
	function.SetInputJsonExample(&model.AdVpn{})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
}

func (s *User) GetVpnEnableList(ctx gtype.Context, ps gtype.Params) {
	ad := s.Ad()
	users, err := ad.GetVpnUsers()
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	results := make(model.AdUserCollection, 0)
	c := len(users)
	for i := 0; i < c; i++ {
		u := users[i]
		if u == nil {
			continue
		}

		results = append(results, &model.AdUser{
			AdDn: model.AdDn{
				Dn: s.ToBase64(u.DN),
			},
			SID:     u.SID,
			Account: u.Account,
			Name:    u.Name,
		})
	}

	sort.Sort(results)
	ctx.Success(results)
}

func (s *User) GetVpnEnableListDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogUser)
	function := catalog.AddFunction(method, uri, "获取VPN用户列表")
	function.SetNote("获取所有VPN状态已启用的帐号列表")
	function.SetOutputDataExample(model.AdUserCollection{
		{
			Name: "张三",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}
