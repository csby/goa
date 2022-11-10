package config

type DhcpApiUri struct {
	Filter DhcpApiUriFilter `json:"filter" note:"筛选器"`
	Lease  DhcpApiUriLease  `json:"lease" note:"地址租用"`
}
