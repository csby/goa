package config

type DhcpApi struct {
	Url string     `json:"url" note:"服务地址"`
	Uri DhcpApiUri `json:"uri" note:"接口地址"`
}
