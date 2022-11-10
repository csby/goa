package user

import (
	"fmt"
	"github.com/csby/goa/assist"
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
	"strings"
)

func NewAccount(log gtype.Log, param *controller.Parameter) *Account {
	instance := &Account{}
	instance.SetLog(log)
	instance.SetParameter(param)

	return instance
}

type Account struct {
	base
}

func (s *Account) CreateUser(ctx gtype.Context, ps gtype.Params) {
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

func (s *Account) CreateUserDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc)
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
