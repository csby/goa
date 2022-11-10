package mail

import (
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
)

func NewBlock(log gtype.Log, param *controller.Parameter) *Block {
	instance := &Block{}
	instance.SetLog(log)
	instance.SetParameter(param)

	return instance
}

type Block struct {
	base
}

func (s *Block) GetReceiverAddressPage(ctx gtype.Context, ps gtype.Params) {
	argument := &gtype.Page{
		PageIndex: 1,
		PageSize:  15,
	}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	uri := uriGetBlockReceiverAddresses
	apiData := &gtype.PageResult{}
	ge := s.callApi(uri, argument, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Block) GetReceiverAddressPageDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "收件人")
	function := catalog.AddFunction(method, uri, "获取阻止列表(分页)")
	function.SetInputJsonExample(&gtype.Page{
		PageIndex: 1,
		PageSize:  15,
	})
	function.SetOutputDataExample(&gtype.PageResult{
		Page: gtype.Page{
			PageIndex: 1,
			PageSize:  15,
		},
		PageCount: 1,
		ItemCount: 101,
		PageItems: []*model.MailBlockAddress{
			{},
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Block) AddReceiverAddress(ctx gtype.Context, ps gtype.Params) {
	argument := &model.MailCommitter{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Address) < 1 {
		ctx.Error(gtype.ErrInput, "地址(address)为空")
		return
	}

	uri := uriAddBlockReceiverAddress
	apiData := false
	ge := s.callApi(uri, argument, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Block) AddReceiverAddressDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "收件人")
	function := catalog.AddFunction(method, uri, "添加阻止地址")
	function.SetInputJsonExample(&model.MailAddress{
		Address: "test@example.com",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Block) DelReceiverAddress(ctx gtype.Context, ps gtype.Params) {
	argument := &model.MailAddress{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Address) < 1 {
		ctx.Error(gtype.ErrInput, "地址(address)为空")
		return
	}

	uri := uriDeleteBlockReceiverAddress
	apiData := false
	ge := s.callApi(uri, argument, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Block) DelReceiverAddressDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "收件人")
	function := catalog.AddFunction(method, uri, "删除阻止地址")
	function.SetInputJsonExample(&model.MailAddress{
		Address: "test@example.com",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Block) GetSenderAddressPage(ctx gtype.Context, ps gtype.Params) {
	argument := &gtype.Page{
		PageIndex: 1,
		PageSize:  15,
	}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	uri := uriGetBlockSenderAddresses
	apiData := &gtype.PageResult{}
	ge := s.callApi(uri, argument, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Block) GetSenderAddressPageDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "发件人")
	function := catalog.AddFunction(method, uri, "获取阻止列表(分页)")
	function.SetInputJsonExample(&gtype.Page{
		PageIndex: 1,
		PageSize:  15,
	})
	function.SetOutputDataExample(&gtype.PageResult{
		Page: gtype.Page{
			PageIndex: 1,
			PageSize:  15,
		},
		PageCount: 1,
		ItemCount: 101,
		PageItems: []*model.MailBlockAddress{
			{},
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Block) AddSenderAddress(ctx gtype.Context, ps gtype.Params) {
	argument := &model.MailCommitter{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Address) < 1 {
		ctx.Error(gtype.ErrInput, "地址(address)为空")
		return
	}

	uri := uriAddBlockSenderAddress
	apiData := false
	ge := s.callApi(uri, argument, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Block) AddSenderAddressDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "发件人")
	function := catalog.AddFunction(method, uri, "添加阻止地址")
	function.SetInputJsonExample(&model.MailAddress{
		Address: "test@example.com",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Block) DelSenderAddress(ctx gtype.Context, ps gtype.Params) {
	argument := &model.MailAddress{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Address) < 1 {
		ctx.Error(gtype.ErrInput, "地址(address)为空")
		return
	}

	uri := uriDeleteBlockSenderAddress
	apiData := false
	ge := s.callApi(uri, argument, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Block) DelSenderAddressDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "发件人")
	function := catalog.AddFunction(method, uri, "删除阻止地址")
	function.SetInputJsonExample(&model.MailAddress{
		Address: "test@example.com",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Block) GetSenderIpPage(ctx gtype.Context, ps gtype.Params) {
	argument := &gtype.Page{
		PageIndex: 1,
		PageSize:  15,
	}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	uri := uriGetRejectSenderIPs
	apiData := &gtype.PageResult{}
	ge := s.callApi(uri, argument, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Block) GetSenderIpPageDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "发送IP")
	function := catalog.AddFunction(method, uri, "获取阻止列表(分页)")
	function.SetInputJsonExample(&gtype.Page{
		PageIndex: 1,
		PageSize:  15,
	})
	function.SetOutputDataExample(&gtype.PageResult{
		Page: gtype.Page{
			PageIndex: 1,
			PageSize:  15,
		},
		PageCount: 1,
		ItemCount: 101,
		PageItems: []*model.MailBlockIP{
			{},
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}
