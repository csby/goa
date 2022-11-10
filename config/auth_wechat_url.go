package config

type AuthWechatUrl struct {
	Api      AuthWechatUrlApi      `json:"api" note:"接口(服务器)配置"`
	Callback AuthWechatUrlCallback `json:"callback" note:"回调地址"`
}
