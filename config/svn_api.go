package config

type SvnApi struct {
	Url string    `json:"url" note:"服务地址"`
	Uri SvnApiUri `json:"uri" note:"接口地址"`
}
