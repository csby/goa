package svn

import (
	"fmt"
	"github.com/csby/goa/controller"
	"github.com/csby/gwsf/gclient"
	"github.com/csby/gwsf/gtype"
)

type base struct {
	controller.Controller
}

func (s *base) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := s.RootCatalog(doc).AddChild("SVN")

	count := len(names)
	if count < 1 {
		return root
	}

	child := root
	for i := 0; i < count; i++ {
		name := names[i]
		child = child.AddChild(name)
	}

	return child
}

func (s *base) callApi(uri string, argument, data interface{}) gtype.Error {
	if s.Cfg == nil {
		return gtype.ErrInternal.SetDetail("cfg is nil")
	}
	baseUrl := s.Cfg.Svn.Api.Url
	if len(baseUrl) < 1 {
		return gtype.ErrInternal.SetDetail("配置错误： 服务地址为空")
	}
	if len(uri) < 1 {
		return gtype.ErrInternal.SetDetail("配置错误： 接口地址为空")
	}
	url := fmt.Sprintf("%s%s", baseUrl, uri)

	client := &gclient.Http{}
	_, output, _, statusCode, err := client.PostJson(url, argument)
	if statusCode != 200 {
		return gtype.ErrInternal.SetDetail("调用接口失败: ", string(output))
	}
	if err != nil {
		return gtype.ErrInternal.SetDetail("调用接口失败: ", err)
	}

	result := &gtype.Result{}
	err = result.Unmarshal(output)
	if err != nil {
		return gtype.ErrInternal.SetDetail("解析接口结果失败: ", err)
	}
	if result.Code != 0 {
		return gtype.NewError(result.Code, result.Error.Summary, nil, result.Error.Detail)
	}

	if data != nil {
		err = result.GetData(data)
		if err != nil {
			return gtype.ErrInternal.SetDetail("解析接口数据失败: ", err)
		}
	}

	return nil
}
