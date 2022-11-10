package auth

import (
	"fmt"
	"github.com/csby/gwsf/gtype"
	"github.com/csby/gwsf/gwechat"
	"strings"
)

type Wechat struct {
	base
}

func (s *Wechat) VerifyUrl(ctx gtype.Context, ps gtype.Params) {
	signature := ""
	timestamp := ""
	nonce := ""
	echostr := ""

	queries := ctx.Queries()
	qc := len(queries)
	for qi := 0; qi < qc; qi++ {
		q := queries[qi]
		if q == nil {
			continue
		}

		qk := strings.ToLower(q.Key)
		if qk == "signature" {
			signature = q.Value
		} else if qk == "timestamp" {
			timestamp = q.Value
		} else if qk == "nonce" {
			nonce = q.Value
		} else if qk == "echostr" {
			echostr = q.Value
		}
	}

	sign := gwechat.SignUrl(s.Cfg.Auth.Wechat.Url.Api.Token, timestamp, nonce)

	w := ctx.Response()
	w.Header().Set("content-type", "text/plain;charset=utf-8")
	if sign != signature {
		fmt.Fprint(w, "")
	} else {
		fmt.Fprint(w, echostr)
	}
}

func (s *Wechat) VerifyUrlDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, authCatalogWechat)
	function := catalog.AddFunction(method, uri, "接口配置")
	function.SetNote("验证微信接口配置信息中的URL相应地址")
	function.SetRemark("验证通过后返回输入的'echostr'")
	function.AddInputQuery(true, "signature", "微信加密签名", "")
	function.AddInputQuery(true, "timestamp", "时间戳", "")
	function.AddInputQuery(true, "nonce", "随机数", "")
	function.AddInputQuery(true, "echostr", "随机字符串", "")
	function.AddOutputHeader("Content-Type", "text/plain; charset=utf-8")
}

func (s *Wechat) GetQRCode(ctx gtype.Context, ps gtype.Params) {
	state := gtype.NewGuid()
	data, err := gwechat.GetLoginPage(s.Cfg.Auth.Wechat.App.ID, s.Cfg.Auth.Wechat.Url.Callback.Code, state)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	ctx.Success(data)
}

func (s *Wechat) GetQRCodeDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, authCatalogWechat)
	function := catalog.AddFunction(method, uri, "获取二维码")
	function.SetNote("获取用于进行扫描的二维码及及登录地址")
	function.SetOutputDataExample(&gwechat.LoginPage{})
}

func (s *Wechat) Code(ctx gtype.Context, ps gtype.Params) {
	code := ""
	state := ""
	queries := ctx.Queries()
	qc := len(queries)
	for qi := 0; qi < qc; qi++ {
		q := queries[qi]
		if q == nil {
			continue
		}

		qk := strings.ToLower(q.Key)
		if qk == "code" {
			code = q.Value
		} else if qk == "state" {
			state = q.Value
		}
	}

	if len(code) > 0 {
		at, ae := gwechat.GetAccessToken(s.Cfg.Auth.Wechat.App.ID, s.Cfg.Auth.Wechat.App.Secret, code)
		if ae == nil && at != nil {
			ui, ue := gwechat.GetUserInfo(at.OpenId, at.AccessToken)
			if ue == nil && ui != nil {
				go s.doAuth(state, ui.OpenId, ui.UnionId, ui.NickName)
			}
		}
	}

	ctx.Response().Write([]byte(wechatHtmlClose))
	ctx.SetHandled(true)
}

func (s *Wechat) CodeDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, authCatalogWechat)
	function := catalog.AddFunction(method, uri, "接收授权码")
	function.SetNote("接受微信服务器回调推送的授权码'")
	function.AddInputQuery(true, "cope", "授权码", "")
	function.AddInputQuery(true, "state", "状态", "")
}

func (s *Wechat) doAuth(state, openId, unionId, nickName string) {
	if len(state) < 1 {
		return
	}

}
