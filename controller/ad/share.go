package ad

import (
	"fmt"
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
	"sort"
	"strings"
)

func NewShare(log gtype.Log, param *controller.Parameter) *Share {
	instance := &Share{}
	instance.SetLog(log)
	instance.SetParameter(param)

	return instance
}

type Share struct {
	base
}

func (s *Share) GetList(ctx gtype.Context, ps gtype.Params) {
	ad := s.Ad()
	root := s.Cfg.Ad.Root.Share
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

func (s *Share) GetListDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogShare)
	function := catalog.AddFunction(method, uri, "获取共享目录列表")
	function.SetNote("获取共享目录列表信息")
	function.SetOutputDataExample([]*model.AdOrganizationUnit{
		{
			Name: "Share",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Share) Add(ctx gtype.Context, ps gtype.Params) {
	argument := &model.AdOrganizationUnitAdd{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput.SetDetail("名称(name)为空"))
		return
	}

	root := s.Cfg.Ad.Root.Share
	if len(root) < 1 {
		ctx.Error(gtype.ErrInternal.SetDetail("配置错误: 根组织单位为空"))
		return
	}
	ad := s.Ad()
	entry, err := ad.GetOrganizationUnit(root)
	if err != nil {
		if ad.IsNotExit(err) {
			entry, err = ad.AddOrganizationUnit(root, "", "")
			if err != nil {
				ctx.Error(gtype.ErrInternal.SetDetail(err))
				return
			}
		} else {
			ctx.Error(gtype.ErrInternal.SetDetail(err))
			return
		}
	}

	ou := fmt.Sprintf("OU=%s,%s", argument.Name, entry.DN)
	entry, err = ad.GetOrganizationUnit(ou)
	if err != nil {
		if !ad.IsNotExit(err) {
			ctx.Error(gtype.ErrInternal.SetDetail(err))
			return
		}
	} else {
		ctx.Error(gtype.ErrInternal.SetDetail(fmt.Sprintf("名称(%s)已存在", argument.Name)))
		return
	}

	entry, err = ad.AddOrganizationUnit(ou, argument.Description, argument.Street)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	groupSuffix := argument.GroupSuffix
	if len(groupSuffix) < 1 {
		items := make([]string, 0)
		names := strings.Split(argument.Name, " ")
		for i := 0; i < len(names); i++ {
			name := strings.TrimSpace(names[i])
			if len(name) > 0 {
				items = append(items, name)
			}
		}

		if len(items) > 0 {
			groupSuffix = strings.Join(items, ".")
		}
	}

	if len(groupSuffix) > 0 {
		groupName := fmt.Sprintf("Share.Authorization.Managers.%s", groupSuffix)
		ad.NewGroup(entry.DN, groupName, "授权管理员", "具备添加、删除成员及编辑成员访问权限的权限")

		groupName = fmt.Sprintf("Share.Read.%s", groupSuffix)
		ad.NewGroup(entry.DN, groupName, "只读用户", "对共享目录具有只读权限")

		groupName = fmt.Sprintf("Share.Read.Write.%s", groupSuffix)
		ad.NewGroup(entry.DN, groupName, "读写用户", "对共享目录具有读写权限，但没有删改权限")

		groupName = fmt.Sprintf("Share.Read.Write.Modify.%s", groupSuffix)
		ad.NewGroup(entry.DN, groupName, "读写删改用户", "对共享目录具有读写及删改权限")
	}

	ctx.Success(nil)
}

func (s *Share) AddDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogShare)
	function := catalog.AddFunction(method, uri, "添加共享目录")
	function.SetInputJsonExample(&model.AdOrganizationUnitAdd{
		Name:        "Public Share",
		GroupSuffix: "Public.Share",
		Description: "总共容量：20.00GB",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}
