package ad

import (
	"fmt"
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
	"sort"
)

func NewServer(log gtype.Log, param *controller.Parameter) *Server {
	instance := &Server{}
	instance.SetLog(log)
	instance.SetParameter(param)

	return instance
}

type Server struct {
	base
}

func (s *Server) GetList(ctx gtype.Context, ps gtype.Params) {
	ad := s.Ad()
	root := s.Cfg.Ad.Root.Server
	units, err := ad.GetOrganizationUnits(root)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	results := make(model.AdServerCollection, 0)
	c := len(units)
	for i := 0; i < c; i++ {
		item := units[i]
		if item == nil {
			continue
		}

		result := &model.AdServer{}
		result.Dn = s.ToBase64(item.DN)
		result.Name = item.Name
		result.Ip = item.Description
		result.Description = item.Street

		results = append(results, result)
	}

	sort.Sort(results)
	ctx.Success(results)
}

func (s *Server) GetListDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogServer)
	function := catalog.AddFunction(method, uri, "获取服务器列表")
	function.SetNote("获取服务器列表信息")
	function.SetOutputDataExample([]*model.AdServer{
		{
			Name: "SVN服务器",
			Ip:   "192.168.1.13",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Server) Add(ctx gtype.Context, ps gtype.Params) {
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

	root := s.Cfg.Ad.Root.Server
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
		groupSuffix = s.CreateAdler32String(argument.Name)
	}

	if len(groupSuffix) > 0 {
		groupName := fmt.Sprintf("Server.Authorization.Managers.%s", groupSuffix)
		ad.NewGroup(entry.DN, groupName, "授权管理员", "具备添加、删除成员及编辑成员访问权限的权限")

		groupName = fmt.Sprintf("Server.Remote.Desktop.Users.%s", groupSuffix)
		ad.NewGroup(entry.DN, groupName, "远程桌面用户", "允许通过远程桌面服务登陆")
	}

	ctx.Success(nil)
}

func (s *Server) AddDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, adCatalogServer)
	function := catalog.AddFunction(method, uri, "添加服务器")
	function.SetInputJsonExample(&model.AdOrganizationUnitAdd{
		Name:        "即时通讯服务器",
		GroupSuffix: "IM",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}
