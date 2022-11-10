package assist

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strings"
)

func (s *Ad) Login(account, password string) (*AdEntryUser, error) {
	if len(account) < 1 {
		return nil, fmt.Errorf("帐号为空")
	}
	samAccount := s.toSamAccount(account)
	if len(samAccount) < 1 {
		return nil, fmt.Errorf("帐号(%s)无效", account)
	}

	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return s.login(conn, samAccount, password)
}

func (s *Ad) GetAllUsers() (AdEntryUserDict, error) {
	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	searchFilter := fmt.Sprintf("(&(&(objectCategory=%s)(objectClass=%s)))", AdCategoryPerson, AdClassUser)
	searchAttrs := []string{"name", "objectGUID", "objectSid", "sAMAccountName"}
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
	results := make(AdEntryUserDict)
	for _, searchEntry := range searchResult.Entries {
		result := &AdEntryUser{}
		s.copyUser(result, searchEntry)

		results[result.SID] = result
	}

	return results, nil
}

func (s *Ad) GetVpnUsers() ([]*AdEntryUser, error) {
	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	filter := &AdEntryFilter{}
	filter.Dialing = "TRUE"

	return s.getUsers(conn, filter)
}

func (s *Ad) GetUsers(parentDN string) ([]*AdEntryUser, error) {
	if len(parentDN) < 1 {
		return nil, fmt.Errorf("parent distinguished name is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	filter := &AdEntryFilter{}
	filter.ParentDN = parentDN

	return s.getUsers(conn, filter)
}

func (s *Ad) GetUserSubordinates(account string) ([]*AdEntryUser, error) {
	if len(account) < 1 {
		return nil, fmt.Errorf("account is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	filter := &AdEntryFilter{}
	filter.Account = account

	users, err := s.getUsers(conn, filter)
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, fmt.Errorf("帐号(%s)不存在", account)
	}
	user := users[0]
	if user == nil {
		return nil, fmt.Errorf("帐号(%s)无效", account)
	}

	return s.getUsers(conn, &AdEntryFilter{Manager: user.DN})
}

func (s *Ad) GetUsersFromGroup(groupDN string) ([]*AdEntryUser, error) {
	if len(groupDN) < 1 {
		return nil, fmt.Errorf("group name is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return s.getUsersFromGroup(conn, groupDN)
}

func (s *Ad) NewUser(v *AdEntryUserCreate) (*AdEntryUser, error) {
	if v == nil {
		return nil, fmt.Errorf("parameter is nil")
	}

	if len(v.Name) < 0 {
		return nil, fmt.Errorf("用户姓名为空")
	}
	if len(v.Account) < 0 {
		return nil, fmt.Errorf("登录帐号为空")
	}
	if len(v.Password) < 0 {
		return nil, fmt.Errorf("登录密码为空")
	}
	pwd, err := s.encodePassword(v.Password)
	if err != nil {
		return nil, fmt.Errorf("登录密码无效: %v", err)
	}

	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	userDn := fmt.Sprintf("CN=%s,CN=Users,%s", v.Name, s.Base)
	if len(v.Parent) > 0 {
		parent, err := s.getOrganizationUnit(conn, &AdEntryFilter{DNs: []string{v.Parent}})
		if err != nil {
			if s.IsNotExit(err) {
				return nil, fmt.Errorf("组织单位(%s)不存在", s.GetDnName(v.Parent))
			} else {
				return nil, err
			}
		}

		userDn = fmt.Sprintf("CN=%s,%s", v.Name, parent.DN)
	}

	if len(v.Manager) > 0 {
		users, err := s.getUsers(conn, &AdEntryFilter{DNs: []string{v.Manager}})
		if err != nil {
			return nil, err
		}
		if len(users) < 0 {
			return nil, fmt.Errorf("直接主管(%s)不存在", s.GetDnName(v.Parent))
		}
	}

	users, err := s.getUsers(conn, &AdEntryFilter{DNs: []string{userDn}})
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return nil, fmt.Errorf("用户姓名(%s)已存在", v.Name)
	}
	users, err = s.getUsers(conn, &AdEntryFilter{Account: v.Account})
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return nil, fmt.Errorf("登录帐号(%s)已存在", v.Account)
	}

	addRequest := ldap.NewAddRequest(userDn, nil)
	addRequest.Attribute("objectClass", []string{AdClassUser})
	addRequest.Attribute("sAMAccountName", []string{v.Account})
	if len(v.Manager) > 0 {
		addRequest.Attribute("manager", []string{v.Manager})
	}
	err = conn.Add(addRequest)
	if err != nil {
		return nil, err
	}

	modifyRequest := ldap.NewModifyRequest(userDn, nil)
	modifyRequest.Replace("unicodePwd", []string{pwd})
	err = conn.Modify(modifyRequest)
	if err != nil {
		s.deleteEntry(conn, userDn)
		le, ok := err.(*ldap.Error)
		if ok {
			if le.ResultCode == 53 {
				return nil, fmt.Errorf("登录密码不符合复杂度要求")
			}
		}
		return nil, err
	}

	_, err = s.setUserControl(conn, s.Base, &AdEntryFilter{DNs: []string{userDn}}, &AdEntryUserControl{DontExpirePassword: true})

	users, err = s.getUsers(conn, &AdEntryFilter{DNs: []string{userDn}})
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return users[0], nil
	}

	return nil, nil
}

func (s *Ad) SetUserPassword(account, password string) error {
	if len(account) < 1 {
		return fmt.Errorf("帐号为空")
	}
	samAccount := s.toSamAccount(account)
	if len(samAccount) < 1 {
		return fmt.Errorf("帐号(%s)无效", account)
	}

	conn, err := s.open(true)
	if err != nil {
		return err
	}
	defer conn.Close()

	users, err := s.getUsers(conn, &AdEntryFilter{Account: samAccount})
	if err != nil {
		return err
	}
	if len(users) < 1 {
		return fmt.Errorf("帐号(%s)不存在", account)
	}
	user := users[0]
	if user == nil {
		return fmt.Errorf("帐号(%s)无效", account)
	}

	return s.setUserPassword(conn, user.DN, password)
}

func (s *Ad) ChangeUserPassword(account, oldPassword, newPassword string) error {
	if len(account) < 1 {
		return fmt.Errorf("帐号为空")
	}
	samAccount := s.toSamAccount(account)
	if len(samAccount) < 1 {
		return fmt.Errorf("帐号(%s)无效", account)
	}

	conn, err := s.open(true)
	if err != nil {
		return err
	}
	defer conn.Close()

	user, err := s.login(conn, samAccount, oldPassword)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("帐号(%s)无效", account)
	}

	return s.setUserPassword(conn, user.DN, newPassword)
}

func (s *Ad) SetUserVpnEnable(account string, enable bool) error {
	if len(account) < 1 {
		return fmt.Errorf("帐号为空")
	}
	samAccount := s.toSamAccount(account)
	if len(samAccount) < 1 {
		return fmt.Errorf("帐号(%s)无效", account)
	}

	conn, err := s.open(true)
	if err != nil {
		return err
	}
	defer conn.Close()

	users, err := s.getUsers(conn, &AdEntryFilter{Account: samAccount})
	if err != nil {
		return err
	}
	if len(users) < 1 {
		return fmt.Errorf("帐号(%s)不存在", account)
	}
	user := users[0]
	if user == nil {
		return fmt.Errorf("帐号(%s)无效", account)
	}

	return s.setUserVpnEnable(conn, user.DN, user.Dialing, enable)
}

func (s *Ad) GetUserVpnEnable(account string) (bool, error) {
	if len(account) < 1 {
		return false, fmt.Errorf("帐号为空")
	}
	samAccount := s.toSamAccount(account)
	if len(samAccount) < 1 {
		return false, fmt.Errorf("帐号(%s)无效", account)
	}

	conn, err := s.open(true)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	users, err := s.getUsers(conn, &AdEntryFilter{Account: samAccount})
	if err != nil {
		return false, err
	}
	if len(users) < 1 {
		return false, fmt.Errorf("帐号(%s)不存在", account)
	}
	user := users[0]
	if user == nil {
		return false, fmt.Errorf("帐号(%s)无效", account)
	}

	if strings.ToUpper(user.Dialing) == "TRUE" {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *Ad) login(conn *ldap.Conn, account, password string) (*AdEntryUser, error) {
	if len(account) < 1 {
		return nil, fmt.Errorf("帐号为空")
	}
	samAccount := s.toSamAccount(account)
	if len(samAccount) < 1 {
		return nil, fmt.Errorf("帐号(%s)无效", account)
	}

	filter := &AdEntryFilter{}
	filter.Account = samAccount
	searchFilter := filter.GetFilter(AdClassUser)
	searchAttrs := []string{"name", "objectGUID", "objectSid", "sAMAccountName"}
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
		return nil, s.fmtError(AdErrorNotExist, "帐号(%s)不存在", account)
	}
	searchEntry := searchResult.Entries[0]

	entry := &AdEntryUser{}
	entry.Name = searchEntry.GetAttributeValue("name")
	entry.GUID = s.decodeGUID(searchEntry.GetRawAttributeValue("objectGUID"))
	entry.DN = searchEntry.DN
	entry.SID = s.decodeSID(searchEntry.GetRawAttributeValue("objectSid"))
	entry.Account = searchEntry.GetAttributeValue("sAMAccountName")

	ctrl, ce := s.getUserControl(conn, entry.DN, &AdEntryFilter{DNs: []string{entry.DN}})
	if ce != nil {
		return nil, ce
	}
	if ctrl.Disable {
		return nil, fmt.Errorf("帐号(%s)已禁用", account)
	}

	l, e := s.open(false)
	if e != nil {
		return nil, e
	}
	defer l.Close()
	e = l.Bind(entry.DN, password)
	if e != nil {
		return nil, fmt.Errorf("密码错误")
	}

	return entry, nil
}

func (s *Ad) getUsers(conn *ldap.Conn, filter *AdEntryFilter) ([]*AdEntryUser, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	searchFilter := filter.GetFilter(AdClassUser)
	searchAttrs := []string{"name", "objectGUID", "objectSid", "sAMAccountName", "msNPAllowDialin"}
	base := filter.ParentDN
	if len(base) < 1 {
		base = s.Base
	}
	searchRequest := ldap.NewSearchRequest(
		base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		searchAttrs,
		nil,
	)
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	results := make([]*AdEntryUser, 0)
	for _, searchEntry := range searchResult.Entries {
		result := &AdEntryUser{}
		s.copyUser(result, searchEntry)

		results = append(results, result)
	}

	return results, nil
}

func (s *Ad) getUsersFromGroup(conn *ldap.Conn, groupDN string) ([]*AdEntryUser, error) {
	if len(groupDN) < 1 {
		return nil, fmt.Errorf("group distinguished name is empty")
	}

	filter := &AdEntryFilter{}
	filter.DNs = []string{groupDN}
	searchFilter := filter.GetFilter(AdClassGroup)
	searchAttrs := []string{"member"}
	searchRequest := ldap.NewSearchRequest(
		groupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		searchAttrs,
		nil,
	)
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	uf := &AdEntryFilter{
		DNs: make([]string, 0),
	}

	for _, searchEntry := range searchResult.Entries {
		attributes := searchEntry.GetAttributeValues("member")
		for _, attribute := range attributes {
			if len(attribute) < 1 {
				continue
			}
			uf.DNs = append(uf.DNs, attribute)
		}
	}

	if len(uf.DNs) < 1 {
		return []*AdEntryUser{}, nil
	}

	return s.getUsers(conn, uf)
}

func (s *Ad) getUserControl(conn *ldap.Conn, base string, filter *AdEntryFilter) (*AdEntryUserControl, error) {
	searchFilter := filter.GetFilter(AdClassUser)
	searchAttrs := []string{"userAccountControl"}
	searchRequest := ldap.NewSearchRequest(
		base,
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
		return nil, s.fmtError(AdErrorNotExist, "not found: %s", searchFilter)
	}

	userAccountControl := ""
	for _, item := range searchResult.Entries {
		userAccountControl = item.GetAttributeValue("userAccountControl")
		break
	}

	controlEntry := &AdEntryUserControl{}
	err = controlEntry.FromValue(userAccountControl)
	if err != nil {
		return nil, err
	}

	return controlEntry, nil
}

func (s *Ad) setUserControl(conn *ldap.Conn, base string, filter *AdEntryFilter, control *AdEntryUserControl) (*AdEntryUserControl, error) {
	searchFilter := filter.GetFilter(AdClassUser)
	searchAttrs := []string{"dn", "userAccountControl"}
	searchRequest := ldap.NewSearchRequest(
		base,
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
		return nil, s.fmtError(AdErrorNotExist, "not found: %s", searchFilter)
	}

	userDn := ""
	userAccountControl := ""
	for _, item := range searchResult.Entries {
		userDn = item.DN
		userAccountControl = item.GetAttributeValue("userAccountControl")
		break
	}
	controlValue, err := control.ToValue(userAccountControl)
	if err != nil {
		return nil, err
	}

	modifyRequest := ldap.NewModifyRequest(userDn, nil)
	modifyRequest.Replace("userAccountControl", []string{controlValue})
	err = conn.Modify(modifyRequest)
	if err != nil {
		return nil, err
	}

	return control, nil
}

func (s *Ad) setUserPassword(conn *ldap.Conn, dn, password string) error {
	pwd, err := s.encodePassword(password)
	if err != nil {
		return fmt.Errorf("登录密码无效: %v", err)
	}

	modifyRequest := ldap.NewModifyRequest(dn, nil)
	modifyRequest.Replace("unicodePwd", []string{pwd})
	err = conn.Modify(modifyRequest)
	if err != nil {
		le, ok := err.(*ldap.Error)
		if ok {
			if le.ResultCode == 53 {
				return fmt.Errorf("密码不符合复杂度要求")
			}
		}
		return err
	}

	return nil
}

func (s *Ad) setUserVpnEnable(conn *ldap.Conn, dn, dialing string, enable bool) error {
	modifyRequest := ldap.NewModifyRequest(dn, nil)
	if enable {
		if len(dialing) > 0 {
			modifyRequest.Replace("msNPAllowDialin", []string{"TRUE"})
		} else {
			modifyRequest.Add("msNPAllowDialin", []string{"TRUE"})
		}
	} else {
		if len(dialing) > 0 {
			modifyRequest.Delete("msNPAllowDialin", []string{"TRUE"})
		} else {
			return nil
		}
	}
	err := conn.Modify(modifyRequest)
	if err != nil {
		return err
	}

	return nil
}
