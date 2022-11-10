package config

type Auth struct {
	Token  AuthToken  `json:"token"`
	Wechat AuthWechat `json:"wechat"`
}
