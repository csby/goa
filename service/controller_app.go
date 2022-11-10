package main

import (
	"github.com/csby/goa/controller"
	"github.com/csby/goa/controller/ad"
	"github.com/csby/goa/controller/auth"
	"github.com/csby/goa/controller/dhcp"
	"github.com/csby/goa/controller/mail"
	"github.com/csby/goa/controller/svn"
	"github.com/csby/goa/controller/user"
	"github.com/csby/gwsf/gtype"
)

type controllerApp struct {
	authAd      *auth.Ad
	userLogin   *user.Login
	userNotify  *user.Notify
	userAccount *user.Account

	dhcpFilter *dhcp.Filter

	svnRepository *svn.Repository
	svnGroup      *svn.Group
	svnUser       *svn.User
	svnPermission *svn.Permission

	mailBlock *mail.Block

	adUser   *ad.User
	adGroup  *ad.Group
	adServer *ad.Server
	adShear  *ad.Share
}

func (s *controllerApp) initController(h *Handler) {
	param := &controller.Parameter{}
	param.Cfg = cfg
	param.Tdb = h.tdb
	param.WChs = h.wsc

	s.authAd = auth.NewAd(log, param)
	s.userLogin = user.NewLogin(log, param)
	s.userNotify = user.NewNotify(log, param)
	s.userAccount = user.NewAccount(log, param)
	s.dhcpFilter = dhcp.NewFilter(log, param)
	s.svnRepository = svn.NewRepository(log, param)
	s.svnGroup = svn.NewGroup(log, param)
	s.svnUser = svn.NewUser(log, param)
	s.svnPermission = svn.NewPermission(log, param)
	s.mailBlock = mail.NewBlock(log, param)
	s.adUser = ad.NewUser(log, param)
	s.adGroup = ad.NewGroup(log, param)
	s.adServer = ad.NewServer(log, param)
	s.adShear = ad.NewShare(log, param)
}

func (s *controllerApp) initRouter(router gtype.Router, path *gtype.Path, preHandle gtype.HttpHandle) {
	// 权限管理-AD
	router.POST(path.Uri("/auth/ad/captcha").SetTokenUI(nil).SetTokenCreate(nil), nil,
		s.authAd.GetCaptcha, s.authAd.GetCaptchaDoc)
	router.POST(path.Uri("/auth/ad/login").SetTokenUI(nil).SetTokenCreate(nil), nil,
		s.authAd.Login, s.authAd.LoginDoc)
	router.POST(path.Uri("/auth/ad/logout"), preHandle,
		s.authAd.Logout, s.authAd.LogoutDoc)

	// 用户管理
	router.POST(path.Uri("/user/login/account"), preHandle,
		s.userLogin.GetAccount, s.userLogin.GetAccountDoc)
	router.GET(path.Uri("/user/login/notify").SetTokenPlace(gtype.TokenPlaceQuery).SetIsWebsocket(true), preHandle,
		s.userNotify.Socket, s.userNotify.SocketDoc)
	router.POST(path.Uri("/user/account/create"), preHandle,
		s.userAccount.CreateUser, s.userAccount.CreateUserDoc)

	// DHCP-筛选器
	router.POST(path.Uri("/dhcp/filter/list"), preHandle,
		s.dhcpFilter.GetDhcpFilters, s.dhcpFilter.GetDhcpFiltersDoc)
	router.POST(path.Uri("/dhcp/filter/add"), preHandle,
		s.dhcpFilter.AddDhcpFilter, s.dhcpFilter.AddDhcpFilterDoc)
	router.POST(path.Uri("/dhcp/filter/del"), preHandle,
		s.dhcpFilter.DelDhcpFilter, s.dhcpFilter.DelDhcpFilterDoc)
	router.POST(path.Uri("/dhcp/filter/mod"), preHandle,
		s.dhcpFilter.ModDhcpFilter, s.dhcpFilter.ModDhcpFilterDoc)
	router.POST(path.Uri("/dhcp/lease/list"), preHandle,
		s.dhcpFilter.GetDhcpLeases, s.dhcpFilter.GetDhcpLeasesDoc)

	// SVN
	router.POST(path.Uri("/svn/repository/add"), preHandle,
		s.svnRepository.AddRepository, s.svnRepository.AddRepositoryDoc)
	router.POST(path.Uri("/svn/repository/list"), preHandle,
		s.svnRepository.GetRepositories, s.svnRepository.GetRepositoriesDoc)
	router.POST(path.Uri("/svn/folder/list"), preHandle,
		s.svnRepository.GetFolders, s.svnRepository.GetFoldersDoc)
	router.POST(path.Uri("/svn/group/role/list"), preHandle,
		s.svnGroup.GetGrantGroups, s.svnGroup.GetGrantGroupsDoc)
	//router.POST(path.Uri("/svn/user/all/list"), preHandle,
	//	s.svnUser.GetAll, s.svnUser.GetAllDoc)
	//router.POST(path.Uri("/svn/user/group/list"), preHandle,
	//	s.svnUser.GetGroups, s.svnUser.GetGroupsDoc)
	router.POST(path.Uri("/svn/permission/user/list"), preHandle,
		s.svnUser.GetPermissions, s.svnUser.GetPermissionsDoc)
	router.POST(path.Uri("/svn/permission/item/list"), preHandle,
		s.svnPermission.GetItemList, s.svnPermission.GetItemListDoc)
	router.POST(path.Uri("/svn/permission/item/add"), preHandle,
		s.svnPermission.AddItem, s.svnPermission.AddItemDoc)
	router.POST(path.Uri("/svn/permission/item/mod"), preHandle,
		s.svnPermission.ModItem, s.svnPermission.ModItemDoc)
	router.POST(path.Uri("/svn/permission/item/del"), preHandle,
		s.svnPermission.DelItem, s.svnPermission.DelItemDoc)

	// 邮件-阻止收件人
	router.POST(path.Uri("/mail/receiver/block/address/page"), preHandle,
		s.mailBlock.GetReceiverAddressPage, s.mailBlock.GetReceiverAddressPageDoc)
	router.POST(path.Uri("/mail/receiver/block/address/add"), preHandle,
		s.mailBlock.AddReceiverAddress, s.mailBlock.AddReceiverAddressDoc)
	router.POST(path.Uri("/mail/receiver/block/address/del"), preHandle,
		s.mailBlock.DelReceiverAddress, s.mailBlock.DelReceiverAddressDoc)
	// 邮件-阻止发件人
	router.POST(path.Uri("/mail/sender/block/address/page"), preHandle,
		s.mailBlock.GetSenderAddressPage, s.mailBlock.GetSenderAddressPageDoc)
	router.POST(path.Uri("/mail/sender/block/address/add"), preHandle,
		s.mailBlock.AddSenderAddress, s.mailBlock.AddSenderAddressDoc)
	router.POST(path.Uri("/mail/sender/block/address/del"), preHandle,
		s.mailBlock.DelSenderAddress, s.mailBlock.DelSenderAddressDoc)
	// 邮件-阻止发送IP
	router.POST(path.Uri("/mail/sender/block/ip/page"), preHandle,
		s.mailBlock.GetSenderIpPage, s.mailBlock.GetSenderIpPageDoc)

	// 域控-用户
	router.POST(path.Uri("/ad/user/account/create"), preHandle,
		s.adUser.CreateUser, s.adUser.CreateUserDoc)
	router.POST(path.Uri("/ad/user/account/password/reset"), preHandle,
		s.adUser.ResetPassword, s.adUser.ResetPasswordDoc)
	router.POST(path.Uri("/ad/user/account/password/change"), preHandle,
		s.adUser.ChangePassword, s.adUser.ChangePasswordDoc)
	router.POST(path.Uri("/ad/user/account/list"), preHandle,
		s.adUser.GetAccountList, s.adUser.GetAccountListDoc)
	router.POST(path.Uri("/ad/user/account/tree"), preHandle,
		s.adUser.GetAccountTree, s.adUser.GetAccountTreeDoc)
	router.POST(path.Uri("/ad/user/org/unit/list"), preHandle,
		s.adUser.GetOrganizationUnitList, s.adUser.GetOrganizationUnitListDoc)
	router.POST(path.Uri("/ad/user/subordinate/list"), preHandle,
		s.adUser.GetSubordinates, s.adUser.GetSubordinatesDoc)
	router.POST(path.Uri("/ad/user/vpn/enable/get"), preHandle,
		s.adUser.GetVpnEnable, s.adUser.GetVpnEnableDoc)
	router.POST(path.Uri("/ad/user/vpn/enable/set"), preHandle,
		s.adUser.SetVpnEnable, s.adUser.SetVpnEnableDoc)
	router.POST(path.Uri("/ad/user/vpn/enable/list"), preHandle,
		s.adUser.GetVpnEnableList, s.adUser.GetVpnEnableListDoc)
	// 域控-组
	router.POST(path.Uri("/ad/group/user/list"), preHandle,
		s.adGroup.GetUsers, s.adGroup.GetUsersDoc)
	router.POST(path.Uri("/ad/group/role/list"), preHandle,
		s.adGroup.GetGrantGroups, s.adGroup.GetGrantGroupsDoc)
	router.POST(path.Uri("/ad/group/member/add"), preHandle,
		s.adGroup.AddMember, s.adGroup.AddMemberDoc)
	router.POST(path.Uri("/ad/group/member/remove"), preHandle,
		s.adGroup.RemoveMember, s.adGroup.RemoveMemberDoc)
	// 域控-服务器
	router.POST(path.Uri("/ad/server/list"), preHandle,
		s.adServer.GetList, s.adServer.GetListDoc)
	router.POST(path.Uri("/ad/server/add"), preHandle,
		s.adServer.Add, s.adServer.AddDoc)
	// 域控-共享目录
	router.POST(path.Uri("/ad/share/list"), preHandle,
		s.adShear.GetList, s.adShear.GetListDoc)
	router.POST(path.Uri("/ad/share/add"), preHandle,
		s.adShear.Add, s.adShear.AddDoc)
}

func (s *controllerApp) createTokenForAccountPassword() func(items []gtype.TokenAuth, ctx gtype.Context) (string, gtype.Error) {
	if s.authAd == nil {
		return nil
	}

	return s.authAd.CreateTokenForAccountPassword
}

func (s *controllerApp) checkToken() func(ctx gtype.Context, ps gtype.Params) {
	if s.authAd == nil {
		return nil
	}

	return s.authAd.CheckToken
}
