package model

import "strings"

type SvnRepository struct {
	Id   string `json:"id" note:"标识ID"`
	Name string `json:"name" required:"true" note:"名称"`
	Path string `json:"path" required:"true" note:"路径"`
}

type SvnRepositoryItem struct {
	Id         string `json:"id" note:"标识ID"`
	Repository string `json:"repository" note:"存储库名称"`
	Name       string `json:"name" note:"项目名称"`
	Path       string `json:"path" note:"路径"`
	Type       int    `json:"type" note:"类型: 0-存储库; 1-文件夹; 2-文件"`
	Url        string `json:"url" note:"地址"`
	Revisions  int    `json:"revisions" note:"修订次数"`

	Children []*SvnRepositoryItem `json:"children" note:"子项"`
}

type SvnRepositoryCreate struct {
	Name string `json:"name" required:"true" note:"名称"`
}

type SvnPermissionID struct {
	AccountId string `json:"accountId" note:"账号ID"`
}

type SvnPermission struct {
	AccountId   string `json:"accountId" note:"账号ID"`
	AccountName string `json:"accountName" note:"账号名称"`
	AccessLevel int    `json:"accessLevel" note:"访问权限: 0-无; 1-只读; 2-读写"`
	Inherited   bool   `json:"inherited" note:"是否继承"`
}

type SvnPermissionArgument struct {
	Repository string `json:"repository" required:"true" note:"存储库名称"`
	Path       string `json:"path" required:"true" note:"路径"`
	AccountId  string `json:"accountId" required:"true" note:"账号ID"`
}

type SvnPermissionArgumentEdit struct {
	SvnPermissionArgument
	AccessLevel int `json:"accessLevel" note:"访问权限: 0-无; 1-只读; 2-读写"`
}

type SvnUser struct {
	Id      string `json:"id" note:"ID"`
	Name    string `json:"name" note:"名称"`
	Account string `json:"account" note:"帐号"`
}

type SvnUserArray []*SvnUser

func (s SvnUserArray) Len() int      { return len(s) }
func (s SvnUserArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SvnUserArray) Less(i, j int) bool {
	return strings.ToLower(s[i].Account) < strings.ToLower(s[j].Account)
}

type SvnUserGroup struct {
	Name  string `json:"name" note:"名称"`
	Users SvnUserArray
}

type SvnUserPermission struct {
	Repository  string `json:"repository" required:"true" note:"存储库名称"`
	Path        string `json:"path" required:"true" note:"路径"`
	AccessLevel int    `json:"accessLevel" note:"访问权限: 0-无; 1-只读; 2-读写"`
}
