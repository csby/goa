package main

import (
	"fmt"
	"github.com/csby/goa/config"
	"github.com/csby/gwsf/gopt"
	"github.com/csby/gwsf/gtype"
	"net/http"
)

func NewHandler(log gtype.Log, cfg *config.Config) gtype.Handler {
	instance := &Handler{}
	instance.SetLog(log)

	instance.wsc = gtype.NewSocketChannelCollection()
	tokenExpiredMinutes := int64(0)
	if cfg != nil {
		tokenExpiredMinutes = cfg.Site.Opt.Api.Token.Expiration
	}
	instance.tdb = gtype.NewTokenDatabase(tokenExpiredMinutes, "staff")

	return instance
}

type Handler struct {
	gtype.Base

	ctrl controllers

	wsc gtype.SocketChannelCollection
	tdb gtype.TokenDatabase
}

func (s *Handler) InitRouting(router gtype.Router) {
	s.ctrl.auth.initRouter(router, authPath, nil)

	apiPath.DefaultTokenUI = gtype.TokenUIForAccountPassword
	apiPath.DefaultTokenCreate = s.ctrl.app.createTokenForAccountPassword()
	s.ctrl.app.initRouter(router, apiPath, s.ctrl.app.checkToken())
}

func (s *Handler) BeforeRouting(ctx gtype.Context) {
	method := ctx.Method()

	// enable across access
	if method == "OPTIONS" {
		ctx.Response().Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "content-type,token")
		ctx.SetHandled(true)
		return
	}

	// default to app site
	if method == "GET" {
		path := ctx.Path()
		if "/" == path || "" == path || config.StaffSiteUri == path {
			redirectUrl := fmt.Sprintf("%s://%s%s/", ctx.Schema(), ctx.Host(), config.StaffSiteUri)
			http.Redirect(ctx.Response(), ctx.Request(), redirectUrl, http.StatusMovedPermanently)
			ctx.SetHandled(true)
			return
		} else if gopt.WebPath == path {
			redirectUrl := fmt.Sprintf("%s://%s%s/", ctx.Schema(), ctx.Host(), gopt.WebPath)
			http.Redirect(ctx.Response(), ctx.Request(), redirectUrl, http.StatusMovedPermanently)
			ctx.SetHandled(true)
			return
		}
	}
}

func (s *Handler) AfterRouting(ctx gtype.Context) {

}

func (s *Handler) ExtendOptSetup(opt gtype.Option) {
	if opt == nil {
		return
	}

	opt.SetCloud(cfg.Cloud.Enabled)
	opt.SetNode(cfg.Node.Enabled)
}

func (s *Handler) ExtendOptApi(router gtype.Router,
	path *gtype.Path,
	preHandle gtype.HttpHandle,
	wsc gtype.SocketChannelCollection,
	tdb gtype.TokenDatabase) {
	s.initController()
	s.ctrl.opt.initRouter(router, path, preHandle)
}
