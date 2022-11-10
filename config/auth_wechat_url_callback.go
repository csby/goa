package config

type AuthWechatUrlCallback struct {
	Code string `json:"code" note:"用户授权回调地址"`
}

func (s *AuthWechatUrlCallback) CopyTo(target *AuthWechatUrlCallback) {
	if target == nil {
		return
	}

	target.Code = s.Code
}
