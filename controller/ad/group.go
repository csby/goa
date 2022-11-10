package ad

import (
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
	"sort"
)

func NewGroup(log gtype.Log, param *controller.Parameter) *Group {
	instance := &Group{}
	instance.SetLog(log)
	instance.SetParameter(param)

	return instance
}

type Group struct {
	base
}

func (s *Group) GetUsers(ctx gtype.Context, ps gtype.Params) {
	argument := &model.AdDn{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Dn) < 1 {
		ctx.Error(gtype.ErrInput.SetDetail("dn为空"))
		return
	}
	dn, err := s.FromBase64(argument.Dn)
	if err != nil {
		ctx.Error(gtype.ErrInput, "dn不是有效base64字符: ", err)
		return
	}

	ad := s.Ad()
	items, err := ad.GetUsersFromGroup(dn)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	results := make(model.AdUserCollection, 0)
	c := len(items)
	for i := 0; i < c; i++ {
		item := items[i]
		if item == nil {
			continue
		}

		result := &model.AdUser{}
		result.Dn = s.ToBase64(item.DN)
		result.Account = item.Account
		result.Name = item.Name

		results = append(results, result)
	}

	sort.Sort(results)
	ctx.Success(results)
}

func (s *Group) GetUsersDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogGroup)
	function := catalog.AddFunction(method, uri, "获取用户列表")
	function.SetInputJsonExample(&model.AdDn{})
	function.SetOutputDataExample([]*model.AdUser{
		{
			Account: "admin",
			Name:    "管理员",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Group) GetGrantGroups(ctx gtype.Context, ps gtype.Params) {
	argument := &model.AdDn{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Dn) < 1 {
		ctx.Error(gtype.ErrInput.SetDetail("dn为空"))
		return
	}
	dn, err := s.FromBase64(argument.Dn)
	if err != nil {
		ctx.Error(gtype.ErrInput, "dn不是有效base64字符: ", err)
		return
	}

	ad := s.Ad()
	items, err := ad.GetGroupsFromOrganizationUnit(dn)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	results := make(model.AdRoleGroupCollection, 0)
	c := len(items)
	for i := 0; i < c; i++ {
		item := items[i]
		if item == nil {
			continue
		}

		result := &model.AdRoleGroup{}
		result.Dn = s.ToBase64(item.DN)
		result.Account = item.Account
		result.Description = item.Description
		result.Info = item.Info
		result.Role = s.GetAdGroupRole(item.Account)

		results = append(results, result)
	}

	sort.Sort(results)
	ctx.Success(results)
}

func (s *Group) GetGrantGroupsDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogGroup)
	function := catalog.AddFunction(method, uri, "获取角色组列表")
	function.SetInputJsonExample(&model.AdDn{})
	function.SetOutputDataExample([]*model.AdRoleGroup{
		{
			Role: 0,
		},
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Group) AddMember(ctx gtype.Context, ps gtype.Params) {
	argument := &model.AdGroupMemberArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.GroupDn) < 1 {
		ctx.Error(gtype.ErrInput.SetDetail("groupDn为空"))
		return
	}
	groupDn, err := s.FromBase64(argument.GroupDn)
	if err != nil {
		ctx.Error(gtype.ErrInput, "groupDn不是有效base64字符: ", err)
		return
	}

	if len(argument.MemberDn) < 1 {
		ctx.Error(gtype.ErrInput.SetDetail("memberDn为空"))
		return
	}
	memberDn, err := s.FromBase64(argument.MemberDn)
	if err != nil {
		ctx.Error(gtype.ErrInput, "memberDn不是有效base64字符: ", err)
		return
	}

	ad := s.Ad()
	err = ad.AddGroupMember(groupDn, memberDn)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	ctx.Success(nil)
}

func (s *Group) AddMemberDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogGroup)
	function := catalog.AddFunction(method, uri, "添加组成员")
	function.SetInputJsonExample(&model.AdGroupMemberArgument{})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Group) RemoveMember(ctx gtype.Context, ps gtype.Params) {
	argument := &model.AdGroupMemberArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.GroupDn) < 1 {
		ctx.Error(gtype.ErrInput.SetDetail("groupDn为空"))
		return
	}
	groupDn, err := s.FromBase64(argument.GroupDn)
	if err != nil {
		ctx.Error(gtype.ErrInput, "groupDn不是有效base64字符: ", err)
		return
	}

	if len(argument.MemberDn) < 1 {
		ctx.Error(gtype.ErrInput.SetDetail("memberDn为空"))
		return
	}
	memberDn, err := s.FromBase64(argument.MemberDn)
	if err != nil {
		ctx.Error(gtype.ErrInput, "memberDn不是有效base64字符: ", err)
		return
	}

	ad := s.Ad()
	err = ad.RemoveGroupMember(groupDn, memberDn)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	ctx.Success(nil)
}

func (s *Group) RemoveMemberDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogGroup)
	function := catalog.AddFunction(method, uri, "移除组成员")
	function.SetInputJsonExample(&model.AdGroupMemberArgument{})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}
