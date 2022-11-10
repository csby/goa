package config

import "github.com/csby/gwsf/gcfg"

const (
	StaffSiteUri  = "/staff"
	StaffSiteName = "办公系统"
)

type OA struct {
	Token gcfg.Token `json:"token"`
}
