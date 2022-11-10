package svn

import (
	"fmt"
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

func (s *Group) GetGrantGroups(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnRepositoryCreate{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput.SetDetail("存储库名称(name)为空"))
		return
	}

	root := s.Cfg.Ad.Root.Svn
	if len(root) < 1 {
		ctx.Error(gtype.ErrInternal.SetDetail("配置错误: 根组织单位为空"))
		return
	}

	dn := fmt.Sprintf("OU=%s,%s", argument.Name, root)
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
	catalog := s.createCatalog(doc, "组")
	function := catalog.AddFunction(method, uri, "获取角色组列表")
	function.SetInputJsonExample(&model.SvnRepositoryCreate{})
	function.SetOutputDataExample([]*model.AdRoleGroup{
		{
			Role: 0,
		},
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}
