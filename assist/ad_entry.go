package assist

import (
	"fmt"
	"strconv"
	"strings"
)

type AdEntry struct {
	Name string // name
	GUID string // objectGUID
	DN   string // distinguishedName
}

type AdEntryFilter struct {
	Account  string   // sAMAccountName
	GUID     string   // objectGUID
	SID      string   // objectSid
	DNs      []string // distinguishedName
	ParentDN string   // msDS-parentdistname
	Manager  string   // manager
	Dialing  string   // msNPAllowDialin
}

func (s *AdEntryFilter) GetFilter(objectClass string) string {
	sb := &strings.Builder{}
	if len(objectClass) > 0 {
		sb.WriteString(fmt.Sprintf("(objectClass=%s)", objectClass))
	}
	if len(s.ParentDN) > 0 {
		sb.WriteString(fmt.Sprintf("(msDS-parentdistname=%s)", s.toFilterValue(s.ParentDN)))
	}
	if len(s.Manager) > 0 {
		sb.WriteString(fmt.Sprintf("(manager=%s)", s.toFilterValue(s.Manager)))
	}
	if len(s.Dialing) > 0 {
		sb.WriteString(fmt.Sprintf("(msNPAllowDialin=%s)", s.toFilterValue(s.Dialing)))
	}
	if len(s.SID) > 0 {
		sb.WriteString(fmt.Sprintf("(objectSid=%s)", s.SID))
	}
	if len(s.GUID) > 0 {
		sb.WriteString(fmt.Sprintf("(objectGUID=%s)", s.GUID))
	}
	if len(s.DNs) > 0 {
		dns := &strings.Builder{}
		dnc := 0
		for _, n := range s.DNs {
			if len(n) < 1 {
				continue
			}
			dns.WriteString(fmt.Sprintf("(distinguishedName=%s)", s.toFilterValue(n)))
			dnc++
		}
		if dnc > 1 {
			sb.WriteString(fmt.Sprintf("(|%s)", dns.String()))
		} else {
			sb.WriteString(dns.String())
		}
	}
	if len(s.Account) > 0 {
		sb.WriteString(fmt.Sprintf("(sAMAccountName=%s)", s.toFilterValue(s.Account)))
	}

	return fmt.Sprintf("(&%s)", sb.String())
}

func (s *AdEntryFilter) toFilterValue(v string) string {
	/*
		*   -> 2a
		(   -> 28
		)   -> 29
		\   -> 5c
		NUL -> 00
	*/
	v = strings.ReplaceAll(v, "*", "\\2a")
	v = strings.ReplaceAll(v, "(", "\\28")
	v = strings.ReplaceAll(v, ")", "\\29")
	v = strings.ReplaceAll(v, "\\", "\\5c")
	v = strings.ReplaceAll(v, "NUL", "\\00")

	return v
}

type AdEntryUserControl struct {
	Disable            bool // 账户已禁用
	DontExpirePassword bool // 密码永不过期
}

func (s *AdEntryUserControl) ToValue(value string) (string, error) {
	val, err := strconv.Atoi(value)
	if err != nil {
		return "", err
	}

	if s.Disable {
		val |= AdAccountDisable
	} else {
		val &= ^AdAccountDisable
	}

	if s.DontExpirePassword {
		val |= AdDontExpirePasswd
	} else {
		val &= ^AdDontExpirePasswd
	}

	return fmt.Sprint(val), nil
}

func (s *AdEntryUserControl) FromValue(value string) error {
	val, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	if (val & AdAccountDisable) == 0 {
		s.Disable = false
	} else {
		s.Disable = true
	}

	if (val & AdDontExpirePasswd) == 0 {
		s.DontExpirePassword = false
	} else {
		s.DontExpirePassword = true
	}

	return nil
}

type AdEntryUser struct {
	AdEntry

	SID     string // objectSid
	Account string // sAMAccountName
	Dialing string // msNPAllowDialin
}

// AdEntryUserDict map[sid]*AdEntryUser
type AdEntryUserDict map[string]*AdEntryUser

type AdEntryUserCreate struct {
	Name     string // 用户姓名
	Account  string // 登录帐号
	Password string // 登录密码
	Manager  string // 直接主管DN
	Parent   string // 组织单位DN
}

type AdEntryOrganizationUnit struct {
	AdEntry

	Description string // description
	Street      string // street
}

type AdEntryGroup struct {
	AdEntry

	SID         string // objectSid
	Account     string // sAMAccountName
	Description string // description
	Info        string // info
}
