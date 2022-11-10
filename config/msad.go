package config

type MsAd struct {
	Host       string      `json:"host" note:"主机地址"`
	Port       int         `json:"port" note:"端口, 389或636"`
	Base       string      `json:"base" note:"根路径，如: DC=example,DC=com"`
	Account    MsAdAccount `json:"account" note:"访问帐号帐号"`
	Root       MsAdRoot    `json:"root" note:"根节点"`
	AdminGroup string      `json:"adminGroup" note:"系统管理员组(帐号名称)"`
}
