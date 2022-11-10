package config

type MsAdRoot struct {
	Server string `json:"server" note:"服务器, 如: OU=服务器,DC=example,DC=com"`
	Share  string `json:"share" note:"共享目录, 如: OU=共享目录,DC=example,DC=com"`
	User   string `json:"user" note:"用户帐号, 如: OU=用户账号,DC=example,DC=com"`
	Svn    string `json:"svn" note:"SVN帐号, 如: OU=SVN,DC=example,DC=com"`
}
