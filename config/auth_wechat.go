package config

type AuthWechat struct {
	Url AuthWechatUrl `json:"url" note:"回调地址"`
	App AuthWechatApp `json:"app" note:"开发帐号"`
}
