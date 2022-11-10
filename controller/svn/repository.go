package svn

import (
	"fmt"
	"github.com/csby/goa/assist"
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
)

func NewRepository(log gtype.Log, param *controller.Parameter) *Repository {
	instance := &Repository{}
	instance.SetLog(log)
	instance.SetParameter(param)

	return instance
}

type Repository struct {
	base
}

func (s *Repository) AddRepository(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnRepositoryCreate{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(name)为空")
		return
	}

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Svn.Api.Uri.Repository.Add
	}

	apiData := ""
	ge := s.callApi(uri, argument, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	go s.addManager(apiData)

	ctx.Success(apiData)
}

func (s *Repository) AddRepositoryDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "新建存储库")
	function.SetNote("成功时返回存储库名称, 并创建3个文件夹: branches,tags,trunk")
	function.SetInputJsonExample(&model.SvnRepositoryCreate{
		Name: "test",
	})
	function.SetOutputDataExample("test")
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Repository) GetRepositories(ctx gtype.Context, ps gtype.Params) {
	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Svn.Api.Uri.Repository.List
	}

	apiData := make([]*model.SvnRepositoryItem, 0)
	ge := s.callApi(uri, nil, &apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Repository) GetRepositoriesDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "获取存储库列表")
	function.SetOutputDataExample([]*model.SvnRepositoryItem{
		{
			Id:         gtype.NewGuid(),
			Repository: "test",
			Name:       "test",
			Path:       "/",
			Children:   []*model.SvnRepositoryItem{},
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Repository) GetFolders(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnRepository{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(name)为空")
		return
	}
	if len(argument.Path) < 1 {
		ctx.Error(gtype.ErrInput, "路径(path)为空")
		return
	}

	uri := ""
	if s.Cfg != nil {
		uri = s.Cfg.Svn.Api.Uri.Repository.Folder
	}

	apiData := &model.SvnRepositoryItem{}
	ge := s.callApi(uri, argument, apiData)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	ctx.Success(apiData)
}

func (s *Repository) GetFoldersDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "获取文件夹列表")
	function.SetInputJsonExample(&model.SvnRepository{
		Name: "test",
		Path: "/",
	})
	function.SetOutputDataExample([]*model.SvnRepositoryItem{
		{
			Id: gtype.NewGuid(),
			Children: []*model.SvnRepositoryItem{
				{
					Id:         gtype.NewGuid(),
					Repository: "test",
					Name:       "trunk",
					Path:       "/trunk",
					Children:   []*model.SvnRepositoryItem{},
				},
			},
		},
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Repository) addManager(repository string) {
	if len(repository) < 1 {
		return
	}
	if s.Cfg == nil {
		return
	}
	root := s.Cfg.Ad.Root.Svn
	if len(root) < 1 {
		return
	}

	ad := &assist.Ad{}
	ad.Host = s.Cfg.Ad.Host
	ad.Port = s.Cfg.Ad.Port
	ad.Base = s.Cfg.Ad.Base
	ad.Account = s.Cfg.Ad.Account.Account
	ad.Password = s.Cfg.Ad.Account.Password

	entry, err := ad.GetOrganizationUnit(root)
	if err != nil {
		if ad.IsNotExit(err) {
			entry, err = ad.AddOrganizationUnit(root, "", "")
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	ou := fmt.Sprintf("OU=%s,%s", repository, entry.DN)
	entry, err = ad.GetOrganizationUnit(ou)
	if err != nil {
		if ad.IsNotExit(err) {
			entry, err = ad.AddOrganizationUnit(ou, "", "")
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	groupName := fmt.Sprintf("SVN.Authorization.Managers.%s", repository)
	ad.NewGroup(entry.DN, groupName, "授权管理员", "具备添加、删除成员及编辑成员访问权限的权限")

	groupName = fmt.Sprintf("SVN.Email.Subscribers.%s", repository)
	ad.NewGroup(entry.DN, groupName, "邮件订阅列表", "当SVN变更后，列表中的用户将接收到邮件通知")
}
