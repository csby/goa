package config

type AuthWechatUrlApi struct {
	Token   string `json:"token" note:"令牌"`
	AESKey  string `json:"aesKey" note:"消息加解密密钥"`
	MsgType int    `json:"msgType" note:"消息加解密方式: 0-明文模式; 1-兼容模式; 2-安全模式"`
}

func (s *AuthWechatUrlApi) CopyTo(target *AuthWechatUrlApi) {
	if target == nil {
		return
	}

	target.Token = s.Token
	target.AESKey = s.AESKey
	target.MsgType = s.MsgType
}

type AuthWechatUrlApiInfo struct {
	AuthWechatUrlApi

	Url string `json:"url" note:"验证响应地址"`
}

func (s *AuthWechatUrlApiInfo) CopyFrom(source *AuthWechatUrlApi) {
	if source == nil {
		return
	}

	source.CopyTo(&s.AuthWechatUrlApi)
}
