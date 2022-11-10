package model

import "strings"

const (
	GroupRoleAuthorization   = 99 // 授权管理员
	GroupRoleRemoteDesktop   = 89 // 远程桌面用户
	GroupRoleDatabaseAdmin   = 79 // 数据库实例系统管理员
	GroupRoleReadOnly        = 69 // 只读
	GroupRoleReadWrite       = 68 // 读写
	GroupRoleReadWriteModify = 67 //读写改

	GroupRoleOther = 0
)

type AdDn struct {
	Dn string `json:"dn" required:"true" note:"唯一名称"`
}

type AdAccount struct {
	Account string `json:"account" required:"true" note:"帐号"`
}

type AdVpn struct {
	AdAccount

	Enable bool `json:"enable" note:"是否启用"`
}

type AdSetPassword struct {
	AdAccount

	Password string `json:"password" required:"true" note:"登录密码"`
}

type AdChangePassword struct {
	AdAccount

	OldPassword string `json:"oldPassword" required:"true" note:"原密码"`
	NewPassword string `json:"newPassword" required:"true" note:"新密码"`
}

type AdGroup struct {
	AdDn

	Account     string `json:"account" note:"帐号名称"`
	Description string `json:"description" note:"描述"`
	Info        string `json:"info" note:"注释"`
}

type AdRoleGroup struct {
	AdGroup

	Role int `json:"role" note:"角色: 99-授权管理员; 89-远程桌面用户; 79-数据库实例管理员; 69-只读; 68-读写; 67-读写改; 0-其他"`
}

type AdRoleGroupCollection []*AdRoleGroup

func (s AdRoleGroupCollection) Len() int      { return len(s) }
func (s AdRoleGroupCollection) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s AdRoleGroupCollection) Less(i, j int) bool {
	if s[i].Role > s[j].Role {
		return true
	} else if s[i].Role < s[j].Role {
		return false
	}

	a, _ := Utf8ToGbk(strings.ToLower(s[i].Description))
	b, _ := Utf8ToGbk(strings.ToLower(s[j].Description))
	l := len(b)
	for idx, chr := range a {
		if idx > l-1 {
			return false
		}
		if chr != b[idx] {
			return chr < b[idx]
		}
	}
	return true
}

type AdUser struct {
	AdDn

	SID     string `json:"sid" note:"ID"`
	Account string `json:"account" note:"帐号"`
	Name    string `json:"name" note:"姓名"`
}

type AdUserCollection []*AdUser

func (s AdUserCollection) Len() int      { return len(s) }
func (s AdUserCollection) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s AdUserCollection) Less(i, j int) bool {
	a, _ := Utf8ToGbk(strings.ToLower(s[i].Name))
	b, _ := Utf8ToGbk(strings.ToLower(s[j].Name))
	l := len(b)
	for idx, chr := range a {
		if idx > l-1 {
			return false
		}
		if chr != b[idx] {
			return chr < b[idx]
		}
	}
	return true
}

type AdUserCreate struct {
	Name     string `json:"name" required:"true" note:"用户姓名"`
	Account  string `json:"account" required:"true" note:"登录帐号"`
	Password string `json:"password" required:"true" note:"登录密码"`
	Manager  string `json:"manager" note:"直接主管DN, base64"`
	Parent   string `json:"parent" note:"组织单位DN, base64"`
}

type AdOrganizationUnit struct {
	AdDn

	Name        string           `json:"name" note:"名称"`
	Description string           `json:"description" note:"描述"`
	Street      string           `json:"street" note:"街道"`
	Users       AdUserCollection `json:"users,omitempty" note:"用户"`
}

type AdOrganizationUnitCollection []*AdOrganizationUnit

func (s AdOrganizationUnitCollection) Len() int      { return len(s) }
func (s AdOrganizationUnitCollection) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s AdOrganizationUnitCollection) Less(i, j int) bool {
	a, _ := Utf8ToGbk(strings.ToLower(s[i].Name))
	b, _ := Utf8ToGbk(strings.ToLower(s[j].Name))
	l := len(b)
	for idx, chr := range a {
		if idx > l-1 {
			return false
		}
		if chr != b[idx] {
			return chr < b[idx]
		}
	}
	return true
}

type AdOrganizationUnitAdd struct {
	Name        string `json:"name" required:"true" note:"名称"`
	Description string `json:"description" note:"描述"`
	Street      string `json:"street" note:"街道"`
	GroupSuffix string `json:"groupSuffix" note:"用户组后缀"`
}

type AdServer struct {
	AdDn

	Name        string `json:"name" note:"名称"`
	Ip          string `json:"ip" note:"IP地址"`
	Description string `json:"description" note:"描述"`
}

type AdServerCollection []*AdServer

func (s AdServerCollection) Len() int      { return len(s) }
func (s AdServerCollection) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s AdServerCollection) Less(i, j int) bool {
	a, _ := Utf8ToGbk(strings.ToLower(s[i].Name))
	b, _ := Utf8ToGbk(strings.ToLower(s[j].Name))
	l := len(b)
	for idx, chr := range a {
		if idx > l-1 {
			return false
		}
		if chr != b[idx] {
			return chr < b[idx]
		}
	}
	return true
}

type AdGroupMemberArgument struct {
	GroupDn  string `json:"groupDn" required:"true" note:"组唯一名称"`
	MemberDn string `json:"memberDn" required:"true" note:"成员唯一名称"`
}
