package config

type DhcpApiUriFilter struct {
	List string `json:"list" note:"获取筛选器列表"`
	Add  string `json:"add" note:"添加筛选器"`
	Del  string `json:"del" note:"删除筛选器"`
	Mod  string `json:"mod" note:"修改筛选器"`
}
