package config

type SvnApiUri struct {
	Repository SvnApiUriRepository `json:"repository"`
	Permission SvnApiUriPermission `json:"permission"`
	User       SvnApiUriUser       `json:"user"`
}
