package user

import (
	"github.com/csby/goa/controller"
	"github.com/csby/gwsf/gtype"
	"time"
)

func NewLogin(log gtype.Log, param *controller.Parameter) *Login {
	instance := &Login{}
	instance.SetLog(log)
	instance.SetParameter(param)

	return instance
}

type Login struct {
	base
}

func (s *Login) GetAccount(ctx gtype.Context, ps gtype.Params) {
	token := s.GetToken(ctx.Token())
	if token == nil {
		ctx.Error(gtype.ErrInternal, "凭证无效")
		return
	}

	account := &gtype.LoginAccount{
		Account:   token.UserAccount,
		Name:      token.UserName,
		LoginTime: gtype.DateTime(token.LoginTime),
		LoginIp:   token.LoginIP,
	}
	if s.IsAdmin(account.Account) {
		account.Role = 1
	}
	ctx.Success(account)
}

func (s *Login) GetAccountDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, userCatalogLogin)
	function := catalog.AddFunction(method, uri, "获取账号信息")
	function.SetNote("获取当前登录账号基本信息")
	function.SetOutputDataExample(&gtype.LoginAccount{
		Account:   "admin",
		Name:      "管理员",
		LoginTime: gtype.DateTime(time.Now()),
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
}
