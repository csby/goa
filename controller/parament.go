package controller

import (
	"github.com/csby/goa/config"
	"github.com/csby/gwsf/gtype"
)

type Parameter struct {
	Cfg  *config.Config
	Tdb  gtype.TokenDatabase
	WChs gtype.SocketChannelCollection
}
