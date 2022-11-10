package dhcp

import (
	"fmt"
	"github.com/csby/goa/assist"
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
	"net"
	"sort"
	"strings"
)

func NewFilter(log gtype.Log, param *controller.Parameter) *Filter {
	instance := &Filter{}
	instance.SetLog(log)
	instance.SetParameter(param)

	return instance
}

type Filter struct {
	base
}

func (s *Filter) GetDhcpFilters(ctx gtype.Context, ps gtype.Params) {
	argument := &model.DhcpFilterListArgument{}
	ctx.GetJson(argument)

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Dhcp.Api.Uri.Filter.List
	}
	apiData := make([]*apiFilterItem, 0)
	ge := s.callApi(uri, nil, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	all := &model.DhcpFilterAll{}
	for i := 0; i < len(apiData); i++ {
		item := apiData[i]
		if item == nil {
			continue
		}
		all.Add(item.Address, item.Comment)
	}

	items := all.Items(argument.Owner)
	sort.Sort(items)

	ctx.Success(items)
}

func (s *Filter) GetDhcpFiltersDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "获取筛选器列表")
	function.SetNote("获取IPv4筛选器列表")
	function.SetInputJsonExample(&model.DhcpFilterListArgument{})
	function.SetOutputDataExample([]*model.DhcpFilter{
		{
			Address: "00-1C-23-20-AF-4A",
			Owner:   "张三",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Filter) AddDhcpFilter(ctx gtype.Context, ps gtype.Params) {
	argument := &model.DhcpFilter{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	if len(argument.Address) < 1 {
		ctx.Error(gtype.ErrInput, "MAC地址为空")
		return
	}
	_, err = net.ParseMAC(argument.Address)
	if err != nil {
		ctx.Error(gtype.ErrInput, fmt.Sprintf("MAC地址(%s)无效", argument.Address))
		return
	}
	argument.Address = strings.ToUpper(strings.ReplaceAll(argument.Address, ":", "-"))
	owner := argument.Owner
	if len(argument.Owner) > 0 {
		ov, oe := s.FromBase64(argument.Owner)
		if oe == nil {
			ad := &assist.Ad{}
			owner = ad.GetDnName(ov)
		}
	}

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Dhcp.Api.Uri.Filter.Add
	}
	apiArgument := &apiFilterItem{
		Allow:   true,
		Address: argument.Address,
		Comment: fmt.Sprintf("%s - %s - %s", owner, argument.Type, argument.Remark),
	}
	ge := s.callApi(uri, apiArgument, nil)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(nil)
}

func (s *Filter) AddDhcpFilterDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "添加筛选器")
	function.SetNote("添加IPv4筛选器到允许或拒绝列表")
	function.SetInputJsonExample(&model.DhcpFilter{
		Address: "00-1C-23-20-AF-4A",
		Owner:   "张三",
		Type:    "手机",
		Remark:  "苹果",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Filter) DelDhcpFilter(ctx gtype.Context, ps gtype.Params) {
	argument := &model.DhcpFilterDeleteArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	if len(argument.Address) < 1 {
		ctx.Error(gtype.ErrInput, "MAC地址为空")
		return
	}
	_, err = net.ParseMAC(argument.Address)
	if err != nil {
		ctx.Error(gtype.ErrInput, fmt.Sprintf("MAC地址(%s)无效", argument.Address))
		return
	}
	argument.Address = strings.ToUpper(strings.ReplaceAll(argument.Address, ":", "-"))

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Dhcp.Api.Uri.Filter.Del
	}
	apiArgument := &apiFilterDelete{
		Address: argument.Address,
	}
	ge := s.callApi(uri, apiArgument, nil)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(nil)
}

func (s *Filter) DelDhcpFilterDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "删除筛选器")
	function.SetNote("从IPv4筛选器允许或拒绝列表中删除指定的筛选器")
	function.SetInputJsonExample(&model.DhcpFilterDeleteArgument{
		Address: "00-1C-23-20-AF-4A",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Filter) ModDhcpFilter(ctx gtype.Context, ps gtype.Params) {
	argument := &model.DhcpFilterModifyArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	if len(argument.Address) < 1 {
		ctx.Error(gtype.ErrInput, "原MAC地址为空")
		return
	}
	if len(argument.Filter.Address) < 1 {
		ctx.Error(gtype.ErrInput, "新MAC地址为空")
		return
	}
	_, err = net.ParseMAC(argument.Filter.Address)
	if err != nil {
		ctx.Error(gtype.ErrInput, fmt.Sprintf("MAC地址(%s)无效", argument.Filter.Address))
		return
	}

	oldAddr := strings.ToUpper(strings.ReplaceAll(argument.Address, ":", "-"))
	newAddr := strings.ToUpper(strings.ReplaceAll(argument.Filter.Address, ":", "-"))
	owner := argument.Filter.Owner
	if len(argument.Filter.Owner) > 0 {
		ov, oe := s.FromBase64(argument.Filter.Owner)
		if oe == nil {
			ad := &assist.Ad{}
			owner = ad.GetDnName(ov)
		}
	}

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Dhcp.Api.Uri.Filter.Mod
	}
	apiArgument := &apiFilterModify{
		Address: oldAddr,
		Filter: apiFilterItem{
			Allow:   true,
			Address: newAddr,
			Comment: fmt.Sprintf("%s - %s - %s", owner, argument.Filter.Type, argument.Filter.Remark),
		},
	}
	ge := s.callApi(uri, apiArgument, nil)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(nil)
}

func (s *Filter) ModDhcpFilterDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "修改筛选器")
	function.SetNote("修改IPv4筛选器允许或拒绝列表中已存在的筛选器")
	function.SetInputJsonExample(&model.DhcpFilterModifyArgument{
		Address: "00-1C-23-20-AF-4A",
		Filter: model.DhcpFilter{
			Address: "00-1C-23-20-AF-4B",
			Owner:   "李四",
			Type:    "平板",
			Remark:  "安卓",
		},
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Filter) GetDhcpLeases(ctx gtype.Context, ps gtype.Params) {
	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Dhcp.Api.Uri.Lease.List
	}
	apiData := make([]*model.DhcpLease, 0)
	ge := s.callApi(uri, nil, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Filter) GetDhcpLeasesDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "地址租用")
	function := catalog.AddFunction(method, uri, "获取地址租用列表")
	function.SetNote("获取IPv4地址租用列表")
	function.SetOutputDataExample([]*model.DhcpLease{
		{
			IpV4:    "192.168.1.102",
			Address: "00-1C-23-20-AF-4A",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}
