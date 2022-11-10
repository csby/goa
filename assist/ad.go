package assist

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
	"strings"
)

type Ad struct {
	Host     string
	Port     int
	Base     string
	Account  string
	Password string
}

func (s *Ad) IsExit(err error) bool {
	if err == nil {
		return false
	}

	v, ok := err.(*AdError)
	if !ok {
		return false
	}
	if v.Code == AdErrorExist {
		return true
	}

	return false
}

func (s *Ad) IsNotExit(err error) bool {
	if err == nil {
		return false
	}

	v, ok := err.(*AdError)
	if !ok {
		return false
	}
	if v.Code == AdErrorNotExist {
		return true
	}

	return false
}

func (s *Ad) GetEntry(filter *AdEntryFilter, objectClass string) (*AdEntry, error) {
	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return s.getEntry(conn, filter, objectClass)
}

func (s *Ad) GetDnName(v string) string {
	vs := strings.Split(v, ",")
	if len(vs) < 1 {
		return ""
	}

	ns := strings.Split(vs[0], "=")
	if len(ns) < 2 {
		return ""
	}

	return ns[1]
}

func (s *Ad) GetDnParent(v string) string {
	index := strings.Index(v, ",")
	if index < 0 {
		return ""
	}

	return v[index+1:]
}

func (s *Ad) fmtExistError(format string, a ...interface{}) *AdError {
	return &AdError{
		Code:    AdErrorExist,
		Message: fmt.Sprintf(format, a...),
	}
}

func (s *Ad) fmtError(code int, format string, a ...interface{}) *AdError {
	return &AdError{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}

func (s *Ad) getEntry(conn *ldap.Conn, filter *AdEntryFilter, objectClass string) (*AdEntry, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	searchFilter := filter.GetFilter(objectClass)
	searchAttrs := []string{"name", "objectGUID"}
	searchRequest := ldap.NewSearchRequest(
		s.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		searchAttrs,
		nil,
	)
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	if len(searchResult.Entries) < 1 {
		return nil, fmt.Errorf("not found: %s", searchFilter)
	}
	searchEntry := searchResult.Entries[0]

	entry := &AdEntry{}
	entry.Name = searchEntry.GetAttributeValue("name")
	entry.GUID = s.decodeGUID(searchEntry.GetRawAttributeValue("objectGUID"))
	entry.DN = searchEntry.DN

	return entry, nil
}

func (s *Ad) open(bind bool) (*ldap.Conn, error) {
	var (
		conn *ldap.Conn
		err  error
	)
	server := fmt.Sprintf("%s:%d", s.Host, s.Port)
	if s.Port == 636 {
		conn, err = ldap.DialTLS("tcp", server, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = ldap.Dial("tcp", server)
	}
	if err != nil {
		return nil, err
	}

	if bind {
		err = conn.Bind(s.Account, s.Password)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}

func (s *Ad) deleteEntry(conn *ldap.Conn, dn string) error {
	if len(dn) < 0 {
		return fmt.Errorf("dn is empty")
	}

	delRequest := ldap.NewDelRequest(dn, nil)
	err := conn.Del(delRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *Ad) decodeSID(sid []byte) string {
	if len(sid) < 28 {
		return ""
	}
	strSid := strings.Builder{}
	strSid.WriteString("S-")

	revision := int(sid[0])
	strSid.WriteString(fmt.Sprint(revision))

	countSubAuths := int(sid[1] & 0xFF)
	authority := int(0)
	for i := 2; i <= 7; i++ {
		shift := uint(8 * (5 - (i - 2)))
		authority |= int(sid[i]) << shift
	}
	strSid.WriteString("-")
	strSid.WriteString(fmt.Sprintf("%x", authority))

	offset := 8
	size := 4
	for j := 0; j < countSubAuths; j++ {
		subAuthority := 0
		for k := 0; k < size; k++ {
			subAuthority |= (int(sid[offset+k]) & 0xFF) << uint(8*k)
		}
		strSid.WriteString("-")
		strSid.WriteString(fmt.Sprint(subAuthority))
		offset += size
	}

	return strSid.String()
}

func (s *Ad) decodeGUID(uuid []byte) string {
	if len(uuid) < 16 {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// https://msdn.microsoft.com/en-us/library/cc223248.aspx
func (s *Ad) encodePassword(password string) (string, error) {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return utf16.NewEncoder().String(fmt.Sprintf(`"%s"`, password))
}

func (s *Ad) toSamAccount(account string) string {
	vs := strings.Split(account, "\\")
	c := len(vs)
	if c < 1 {
		return ""
	}

	vs = strings.Split(vs[c-1], "@")
	if len(vs) < 1 {
		return ""
	}

	return vs[0]
}

func (s *Ad) copyGroup(target *AdEntryGroup, source *ldap.Entry) {
	if target == nil || source == nil {
		return
	}

	target.Name = source.GetAttributeValue("name")
	target.GUID = s.decodeGUID(source.GetRawAttributeValue("objectGUID"))
	target.DN = source.DN
	target.SID = s.decodeSID(source.GetRawAttributeValue("objectSid"))
	target.Account = source.GetAttributeValue("sAMAccountName")
	target.Description = source.GetAttributeValue("description")
	target.Info = source.GetAttributeValue("info")
}

func (s *Ad) copyUser(target *AdEntryUser, source *ldap.Entry) {
	if target == nil || source == nil {
		return
	}

	target.Name = source.GetAttributeValue("name")
	target.GUID = s.decodeGUID(source.GetRawAttributeValue("objectGUID"))
	target.DN = source.DN
	target.SID = s.decodeSID(source.GetRawAttributeValue("objectSid"))
	target.Account = source.GetAttributeValue("sAMAccountName")
	target.Dialing = source.GetAttributeValue("msNPAllowDialin")
}
