package config

type SvnApiUriRepository struct {
	Add    string `json:"add" note:"新建存储库"`
	List   string `json:"list" note:"获取存储库列表"`
	Folder string `json:"folder" note:"获取文件夹列表"`
}
