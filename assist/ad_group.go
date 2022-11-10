package assist

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strings"
)

func (s *Ad) NewGroup(parentDN, name, description, info string) (*AdEntryGroup, error) {
	if len(parentDN) < 1 {
		return nil, fmt.Errorf("parent distinguished name is empty")
	}
	if len(name) < 1 {
		return nil, fmt.Errorf("name is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return s.addGroup(conn, parentDN, name, description, info)
}

func (s *Ad) GetGroup(dn string) (*AdEntryGroup, error) {
	if len(dn) < 1 {
		return nil, fmt.Errorf("dn is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return s.getGroup(conn, &AdEntryFilter{DNs: []string{dn}})
}

func (s *Ad) GetGroupsFromOrganizationUnit(ouDN string) ([]*AdEntryGroup, error) {
	if len(ouDN) < 1 {
		return nil, fmt.Errorf("ouDN is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	filter := &AdEntryFilter{}
	filter.ParentDN = ouDN
	searchFilter := filter.GetFilter(AdClassGroup)
	searchAttrs := []string{"name", "objectGUID", "objectSid", "sAMAccountName", "description", "info"}
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

	results := make([]*AdEntryGroup, 0)
	for _, searchEntry := range searchResult.Entries {
		result := &AdEntryGroup{}
		s.copyGroup(result, searchEntry)

		results = append(results, result)
	}

	return results, nil
}

func (s *Ad) AddGroupMember(groupDN, memberDN string) error {
	if len(groupDN) < 1 {
		return fmt.Errorf("group distinguished name is empty")
	}
	if len(memberDN) < 1 {
		return fmt.Errorf("member distinguished name is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return err
	}
	defer conn.Close()

	return s.addGroupMember(conn, groupDN, memberDN)
}

func (s *Ad) RemoveGroupMember(groupDN, memberDN string) error {
	if len(groupDN) < 1 {
		return fmt.Errorf("group distinguished name is empty")
	}
	if len(memberDN) < 1 {
		return fmt.Errorf("member distinguished name is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return err
	}
	defer conn.Close()

	return s.removeGroupMember(conn, groupDN, memberDN)
}

func (s *Ad) IsGroupMember(groupAccount, memberAccount string) (bool, error) {
	if len(groupAccount) < 1 {
		return false, fmt.Errorf("group account is empty")
	}
	if len(memberAccount) < 1 {
		return false, fmt.Errorf("member account is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	return s.isGroupMember(conn, groupAccount, memberAccount)
}

func (s *Ad) addGroup(conn *ldap.Conn, parentDN, name, description, info string) (*AdEntryGroup, error) {
	if len(parentDN) < 1 {
		return nil, fmt.Errorf("parent distinguished name is empty")
	}
	if len(name) < 1 {
		return nil, fmt.Errorf("name is empty")
	}

	dn := fmt.Sprintf("CN=%s,%s", name, parentDN)
	addRequest := ldap.NewAddRequest(dn, nil)
	addRequest.Attribute("objectClass", []string{AdClassGroup})
	addRequest.Attribute("sAMAccountName", []string{name})
	if len(description) > 0 {
		addRequest.Attribute("description", []string{description})
	}
	if len(info) > 0 {
		addRequest.Attribute("info", []string{info})
	}
	err := conn.Add(addRequest)
	if err != nil {
		le, ok := err.(*ldap.Error)
		if ok {
			if le.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				return nil, s.fmtExistError("对象(%s)已存在", dn)
			}
		}
		return nil, err
	}

	return s.getGroup(conn, &AdEntryFilter{DNs: []string{dn}})
}

func (s *Ad) getGroup(conn *ldap.Conn, filter *AdEntryFilter) (*AdEntryGroup, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	searchFilter := filter.GetFilter(AdClassGroup)
	searchAttrs := []string{"name", "objectGUID", "objectSid", "sAMAccountName", "description", "info"}
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
		return nil, ErrNotExist
	}

	searchEntry := searchResult.Entries[0]
	result := &AdEntryGroup{}
	s.copyGroup(result, searchEntry)

	return result, nil
}

func (s *Ad) addGroupMember(conn *ldap.Conn, groupDN, memberDN string) error {
	if len(groupDN) < 1 {
		return fmt.Errorf("group distinguished name is empty")
	}
	if len(memberDN) < 1 {
		return fmt.Errorf("member distinguished name is empty")
	}

	filter := &AdEntryFilter{DNs: []string{groupDN}}
	searchFilter := filter.GetFilter(AdClassGroup)
	searchAttrs := []string{"dn", "member"}
	searchRequest := ldap.NewSearchRequest(
		s.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		searchAttrs,
		nil,
	)
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return err
	}
	if len(searchResult.Entries) < 1 {
		return fmt.Errorf("group distinguished name (%s) not exist", groupDN)
	}

	groupDn := ""
	members := make([]string, 0)
	for _, item := range searchResult.Entries {
		groupDn = item.DN
		members = item.GetAttributeValues("member")
		break
	}
	members = append(members, memberDN)
	modifyRequest := ldap.NewModifyRequest(groupDn, nil)
	modifyRequest.Replace("member", members)

	err = conn.Modify(modifyRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *Ad) removeGroupMember(conn *ldap.Conn, groupDN, memberDN string) error {
	if len(groupDN) < 1 {
		return fmt.Errorf("group distinguished name is empty")
	}
	if len(memberDN) < 1 {
		return fmt.Errorf("member distinguished name is empty")
	}

	filter := &AdEntryFilter{DNs: []string{groupDN}}
	searchFilter := filter.GetFilter(AdClassGroup)
	searchAttrs := []string{"dn", "member"}
	searchRequest := ldap.NewSearchRequest(
		s.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		searchAttrs,
		nil,
	)
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return err
	}
	if len(searchResult.Entries) < 1 {
		return fmt.Errorf("group distinguished name (%s) not exist", groupDN)
	}

	groupDn := ""
	members := make([]string, 0)
	for _, item := range searchResult.Entries {
		groupDn = item.DN
		members = item.GetAttributeValues("member")
		break
	}

	newMembers := make([]string, 0)
	for _, item := range members {
		if strings.ToLower(item) == strings.ToLower(memberDN) {
			continue
		}

		newMembers = append(newMembers, item)
	}

	members = append(members, memberDN)
	modifyRequest := ldap.NewModifyRequest(groupDn, nil)
	modifyRequest.Replace("member", newMembers)

	err = conn.Modify(modifyRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *Ad) isGroupMember(conn *ldap.Conn, groupAccount, memberAccount string) (bool, error) {
	if len(groupAccount) < 1 {
		return false, fmt.Errorf("group account is empty")
	}
	if len(memberAccount) < 1 {
		return false, fmt.Errorf("member account is empty")
	}

	filter := &AdEntryFilter{Account: groupAccount}
	searchFilter := filter.GetFilter(AdClassGroup)
	searchAttrs := []string{"member"}
	searchRequest := ldap.NewSearchRequest(
		s.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		searchAttrs,
		nil,
	)
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return false, err
	}
	if len(searchResult.Entries) < 1 {
		return false, fmt.Errorf("group account (%s) not exist", groupAccount)
	}

	members := make([]string, 0)
	for _, item := range searchResult.Entries {
		members = item.GetAttributeValues("member")
		break
	}

	filter.Account = memberAccount
	user, err := s.getEntry(conn, filter, AdClassUser)
	if err != nil {
		return false, err
	}

	for _, item := range members {
		if strings.ToLower(item) == strings.ToLower(user.DN) {
			return true, nil
		}
	}

	return false, nil
}
