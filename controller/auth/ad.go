package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/socket"
	"github.com/csby/gsecurity/grsa"
	"github.com/csby/gwsf/gtype"
	"github.com/mojocn/base64Captcha"
	"strings"
	"time"
)

const (
	captchaNumberSource       = "1234567890"
	captchaLetterSource       = "ABCDEFGHIJKLMNPQRSTUVWXYZabcedefghijklmnpqrstuvwxyz"
	captchaLetterNumberSource = "ABCDEFGHIJKLMNPQRSTUVWXYZ12345678abcedefghijklmnprstuvwxyz"
)

func NewAd(log gtype.Log, param *controller.Parameter) *Ad {
	instance := &Ad{}
	instance.SetLog(log)
	instance.SetParameter(param)

	instance.errorCount = make(map[string]int)
	instance.captchaStore = base64Captcha.DefaultMemStore
	instance.rsaPrivate.Create(1024)

	return instance
}

type Ad struct {
	base

	errorCount   map[string]int
	captchaStore base64Captcha.Store
	rsaPrivate   grsa.Private

	AccountVerification func(account, password string) gtype.Error
}

func (s *Ad) GetCaptcha(ctx gtype.Context, ps gtype.Params) {
	filter := &gtype.CaptchaFilter{
		Mode:   3,
		Length: 4,
		Width:  100,
		Height: 30,
	}
	err := ctx.GetJson(filter)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	var driver base64Captcha.Driver
	switch filter.Mode {
	case 0:
		driver = base64Captcha.NewDriverString(filter.Height, filter.Width, 0, 0, filter.Length, captchaNumberSource, &filter.BackColor, []string{})
	case 1:
		driver = base64Captcha.NewDriverString(filter.Height, filter.Width, 0, 0, filter.Length, captchaLetterSource, &filter.BackColor, []string{})
	case 2:
		driver = base64Captcha.NewDriverMath(filter.Height, filter.Width, 0, 0, &filter.BackColor, []string{})
	default:
		driver = base64Captcha.NewDriverString(filter.Height, filter.Width, 0, 0, filter.Length, captchaLetterNumberSource, &filter.BackColor, []string{})
	}

	captcha := base64Captcha.NewCaptcha(driver, s.captchaStore)
	captchaId, captchaValue, err := captcha.Generate()
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	data := &gtype.Captcha{
		ID:       captchaId,
		Value:    captchaValue,
		Required: s.captchaRequired(ctx.RIP()),
	}
	rsaPublic, err := s.rsaPrivate.Public()
	if err == nil {
		keyVal, err := rsaPublic.ToMemory()
		if err == nil {
			data.RsaPublicKey = string(keyVal)
		}
	}

	ctx.Success(data)
}

func (s *Ad) GetCaptchaDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, authCatalogAd)
	function := catalog.AddFunction(method, uri, "获取验证码")
	function.SetNote("获取用户登陆需要的验证码信息")
	function.SetRemark("该接口不需要凭证")
	function.SetInputJsonExample(&gtype.CaptchaFilter{
		Mode:   3,
		Length: 4,
		Width:  100,
		Height: 30,
	})

	function.SetOutputDataExample(&gtype.Captcha{
		ID:           "GKSVhVMRAHsyVuXSrMYs",
		Value:        "data:image/png;base64,iVBOR...",
		RsaPublicKey: "-----BEGIN PUBLIC KEY-----...-----END PUBLIC KEY-----",
		Required:     false,
	})
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Ad) Login(ctx gtype.Context, ps gtype.Params) {
	filter := &gtype.LoginFilter{}
	err := ctx.GetJson(filter)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	requireCaptcha := s.captchaRequired(ctx.RIP())
	err = filter.Check(requireCaptcha)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	if requireCaptcha {
		captchaValue := s.captchaStore.Get(filter.CaptchaId, true)
		if strings.ToLower(captchaValue) != strings.ToLower(filter.CaptchaValue) {
			ctx.Error(gtype.ErrLoginCaptchaInvalid)
			return
		}
	}

	pwd := filter.Password
	if strings.ToLower(filter.Encryption) == "rsa" {
		buf, err := base64.StdEncoding.DecodeString(filter.Password)
		if err != nil {
			ctx.Error(gtype.ErrLoginPasswordInvalid, err)
			s.increaseErrorCount(ctx.RIP())
			return
		}

		decryptedPwd, err := s.rsaPrivate.Decrypt(buf)
		if err != nil {
			ctx.Error(gtype.ErrLoginPasswordInvalid, err)
			s.increaseErrorCount(ctx.RIP())
			return
		}
		pwd = string(decryptedPwd)
	}

	login, be, err := s.Authenticate(ctx, filter.Account, pwd)
	if be != nil {
		ctx.Error(be, err)
		s.increaseErrorCount(ctx.RIP())
		return
	}

	ctx.Success(login)
	s.clearErrorCount(ctx.RIP())
}

func (s *Ad) LoginDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, authCatalogAd)
	function := catalog.AddFunction(method, uri, "用户登录")
	function.SetNote("通过用户账号及密码进行登录获取凭证")
	function.SetRemark("连续3次错误将要求输入验证码")
	function.SetInputJsonExample(&gtype.LoginFilter{
		Account:      "admin",
		Password:     "1",
		CaptchaId:    "r4kcmz2E12e0qJQOvqRB",
		CaptchaValue: "1e35",
		Encryption:   "",
	})

	function.SetOutputDataExample(&gtype.Login{
		Token: "71b9b7e2ac6d4166b18f414942ff3481",
	})
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrLoginCaptchaInvalid)
	function.AddOutputError(gtype.ErrLoginAccountNotExit)
	function.AddOutputError(gtype.ErrLoginPasswordInvalid)
	function.AddOutputError(gtype.ErrLoginAccountOrPasswordInvalid)
}

func (s *Ad) Logout(ctx gtype.Context, ps gtype.Params) {
	tv := ctx.Token()
	if len(tv) < 1 {
		ctx.Error(gtype.ErrTokenEmpty)
		return
	}
	_, ok := s.Tdb.Get(tv, false)
	if !ok {
		ctx.Error(gtype.ErrTokenInvalid)
		return
	}

	s.WriteWebSocketMessage(ctx.Token(), socket.WSUserLogout, nil)
	ok = s.Tdb.Del(tv)
	if ok {
	}

	ctx.Success(nil)
}

func (s *Ad) LogoutDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, authCatalogAd)
	function := catalog.AddFunction(method, uri, "退出登录")
	function.SetNote("退出登录, 使当前凭证失效")
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
}

func (s *Ad) CreateTokenForAccountPassword(items []gtype.TokenAuth, ctx gtype.Context) (string, gtype.Error) {
	account := ""
	password := ""
	count := len(items)
	for i := 0; i < count; i++ {
		item := items[i]
		if item.Name == "account" {
			account = item.Value
		} else if item.Name == "password" {
			password = item.Value
		}
	}

	model, code, err := s.Authenticate(ctx, account, password)
	if code != nil {
		return "", code.SetDetail(err)
	}

	if model != nil {
		return model.Token, nil
	}

	return "", nil
}

func (s *Ad) Authenticate(ctx gtype.Context, account, password string) (*gtype.Login, gtype.Error, error) {
	ad := s.Ad()
	user, err := ad.Login(account, password)
	if err != nil {
		return nil, gtype.ErrLoginAccountOrPasswordInvalid, err
	}

	now := time.Now()
	token := &gtype.Token{
		ID:          ctx.NewGuid(),
		UserAccount: user.Account,
		UserName:    user.Name,
		LoginIP:     ctx.RIP(),
		LoginTime:   now,
		ActiveTime:  now,
		Usage:       0,
		Ext:         user,
	}
	s.Tdb.Set(token.ID, token)

	login := &gtype.Login{
		Token:   token.ID,
		Account: token.UserAccount,
		Name:    token.UserName,
	}

	return login, nil, err
}

func (s *Ad) CheckToken(ctx gtype.Context, ps gtype.Params) {
	tokenValue := ctx.Token()
	if len(tokenValue) < 1 {
		ctx.Error(gtype.ErrTokenEmpty)
		ctx.SetHandled(true)
		return
	}

	token, ok := s.Tdb.Get(tokenValue, true)
	if !ok {
		ctx.Error(gtype.ErrTokenInvalid)
		ctx.SetHandled(true)
		return
	}

	tokenModel, ok := token.(*gtype.Token)
	if !ok {
		ctx.Error(gtype.ErrInternal, "类型转换错误(*gtype.Token)")
		ctx.SetHandled(true)
		return
	}

	if tokenModel.LoginIP != ctx.RIP() {
		ctx.Error(gtype.ErrTokenIllegal, fmt.Sprintf("IP不匹配: 当前IP%s, 登录IP%s", ctx.RIP(), tokenModel.LoginIP))
		ctx.SetHandled(true)
		return
	}
}

func (s *Ad) onWebsocketWriteFilter(message *gtype.SocketMessage, channel gtype.SocketChannel, token *gtype.Token) bool {
	if message == nil {
		return false
	}

	if channel == nil {
		return false
	}

	if token == nil {
		return false
	}
	channelToken := channel.Token()
	if channelToken == nil {
		return false
	}

	return false
}

func (s *Ad) captchaRequired(ip string) bool {
	if s.errorCount == nil {
		return false
	}

	count, ok := s.errorCount[ip]
	if ok {
		if count < 3 {
			return false
		} else {
			return true
		}
	}

	return false
}

func (s *Ad) increaseErrorCount(ip string) {
	if s.errorCount == nil {
		return
	}

	count := 1
	v, ok := s.errorCount[ip]
	if ok {
		count += v
	}

	s.errorCount[ip] = count
}

func (s *Ad) clearErrorCount(ip string) {
	if s.errorCount == nil {
		return
	}

	_, ok := s.errorCount[ip]
	if ok {
		delete(s.errorCount, ip)
	}
}
