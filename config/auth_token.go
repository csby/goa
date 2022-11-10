package config

import (
	"github.com/csby/gwsf/gcfg"
)

type AuthToken struct {
	Code    gcfg.Token `json:"code"`
	Access  gcfg.Token `json:"access"`
	Refresh gcfg.Token `json:"refresh"`
}
