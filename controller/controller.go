package controller

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/csby/goa/assist"
	"github.com/csby/goa/config"
	"github.com/csby/goa/data/model"
	"github.com/csby/gwsf/gtype"
	"hash/adler32"
	"strings"
)

type Controller struct {
	gtype.Base

	Cfg  *config.Config
	Tdb  gtype.TokenDatabase
	WChs gtype.SocketChannelCollection
}

func (s *Controller) SetParameter(p *Parameter) {
	if p == nil {
		return
	}

	s.Cfg = p.Cfg
	s.Tdb = p.Tdb
	s.WChs = p.WChs
}

func (s *Controller) RootCatalog(doc gtype.Doc) gtype.Catalog {
	return doc.AddCatalog("OA平台接口")
}

func (s *Controller) ToBase64(v string) string {
	return base64.URLEncoding.EncodeToString([]byte(v))
}

func (s *Controller) FromBase64(v string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(v)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *Controller) GetToken(key string) *gtype.Token {
	if len(key) < 1 {
		return nil
	}

	if s.Tdb == nil {
		return nil
	}

	value, ok := s.Tdb.Get(key, false)
	if !ok {
		return nil
	}

	token, ok := value.(*gtype.Token)
	if !ok {
		return nil
	}

	return token
}

func (s *Controller) WriteWebSocketMessage(token string, id int, data interface{}) bool {
	if s.WChs == nil {
		return false
	}

	msg := &gtype.SocketMessage{
		ID:   id,
		Data: data,
	}

	s.WChs.Write(msg, s.GetToken(token))

	return true
}

func (s *Controller) CreateAdler32String(a ...interface{}) string {
	h := adler32.New()
	_, err := h.Write([]byte(fmt.Sprint(a...)))
	if err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}

func (s *Controller) Ad() *assist.Ad {
	ad := &assist.Ad{}

	cfg := s.Cfg
	if cfg != nil {
		ad.Host = cfg.Ad.Host
		ad.Port = cfg.Ad.Port
		ad.Base = cfg.Ad.Base
		ad.Account = cfg.Ad.Account.Account
		ad.Password = cfg.Ad.Account.Password
	}

	return ad
}

func (s *Controller) GetAdGroupRole(account string) int {
	av := strings.ToLower(account)
	if strings.Contains(av, ".authorization.") {
		return model.GroupRoleAuthorization
	} else if strings.Contains(av, ".remote.desktop.") {
		return model.GroupRoleRemoteDesktop
	} else if strings.Contains(av, ".database.sysadmin.") {
		return model.GroupRoleDatabaseAdmin
	} else if strings.Contains(av, ".read.write.modify.") {
		return model.GroupRoleReadWriteModify
	} else if strings.Contains(av, ".read.write.") {
		return model.GroupRoleReadWrite
	} else if strings.Contains(av, ".read.") {
		return model.GroupRoleReadOnly
	}

	return model.GroupRoleOther
}

func (s *Controller) IsAdmin(account string) bool {
	if s.Cfg == nil {
		return false
	}

	ad := s.Ad()
	ok, err := ad.IsGroupMember(s.Cfg.Ad.AdminGroup, account)
	if err != nil {
		return false
	}

	return ok
}
