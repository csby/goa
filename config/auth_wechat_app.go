package config

type AuthWechatApp struct {
	ID     string `json:"id" note:"开发帐号标识ID"`
	Secret string `json:"secret" note:"开发帐号密钥"`
}

func (s *AuthWechatApp) CopyTo(target *AuthWechatApp) {
	if target == nil {
		return
	}

	target.ID = s.ID
	target.Secret = s.Secret
}
