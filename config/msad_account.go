package config

type MsAdAccount struct {
	Account  string `json:"account" note:"账号，全路径，如: CN=Administrator,CN=Users,DC=example,DC=com"`
	Password string `json:"password" note:"账号密码"`
}
