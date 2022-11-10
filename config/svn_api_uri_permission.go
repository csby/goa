package config

type SvnApiUriPermission struct {
	List string `json:"list" note:"获取项目访问权限列表"`
	Add  string `json:"add" note:"添加访问权限"`
	Mod  string `json:"mod" note:"修改访问权限"`
	Del  string `json:"del" note:"删除访问权限"`
}
