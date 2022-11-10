package assist

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

func (s *Ad) GetOrganizationUnits(parentDN string) ([]*AdEntryOrganizationUnit, error) {
	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return s.getOrganizationUnits(conn, parentDN)
}

func (s *Ad) AddOrganizationUnit(dn, description, street string) (*AdEntryOrganizationUnit, error) {
	if len(dn) < 1 {
		return nil, fmt.Errorf("distinguished name is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return s.addOrganizationUnit(conn, dn, description, street)
}

func (s *Ad) GetOrganizationUnit(dn string) (*AdEntryOrganizationUnit, error) {
	if len(dn) < 1 {
		return nil, fmt.Errorf("dn is empty")
	}

	conn, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return s.getOrganizationUnit(conn, &AdEntryFilter{DNs: []string{dn}})
}

func (s *Ad) getOrganizationUnits(conn *ldap.Conn, parentDN string) ([]*AdEntryOrganizationUnit, error) {
	if len(parentDN) < 1 {
		return nil, fmt.Errorf("parent distinguished name is empty")
	}

	filter := &AdEntryFilter{}
	filter.ParentDN = parentDN
	searchFilter := filter.GetFilter(AdClassOrganizationalUnit)
	searchAttrs := []string{"name", "objectGUID", "description", "street"}
	searchRequest := ldap.NewSearchRequest(
		parentDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		searchAttrs,
		nil,
	)
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	results := make([]*AdEntryOrganizationUnit, 0)
	for _, searchEntry := range searchResult.Entries {
		result := &AdEntryOrganizationUnit{}
		result.Name = searchEntry.GetAttributeValue("name")
		result.GUID = s.decodeGUID(searchEntry.GetRawAttributeValue("objectGUID"))
		result.DN = searchEntry.DN
		result.Description = searchEntry.GetAttributeValue("description")
		result.Street = searchEntry.GetAttributeValue("street")

		results = append(results, result)
	}

	return results, nil
}

func (s *Ad) getOrganizationUnit(conn *ldap.Conn, filter *AdEntryFilter) (*AdEntryOrganizationUnit, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	searchFilter := filter.GetFilter(AdClassOrganizationalUnit)
	searchAttrs := []string{"name", "objectGUID", "description", "street"}
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
	result := &AdEntryOrganizationUnit{}
	result.Name = searchEntry.GetAttributeValue("name")
	result.GUID = s.decodeGUID(searchEntry.GetRawAttributeValue("objectGUID"))
	result.DN = searchEntry.DN
	result.Description = searchEntry.GetAttributeValue("description")
	result.Street = searchEntry.GetAttributeValue("street")

	return result, nil
}

func (s *Ad) addOrganizationUnit(conn *ldap.Conn, dn, description, street string) (*AdEntryOrganizationUnit, error) {
	if len(dn) < 1 {
		return nil, fmt.Errorf("distinguished name is empty")
	}

	addRequest := ldap.NewAddRequest(dn, nil)
	addRequest.Attribute("objectClass", []string{AdClassOrganizationalUnit})
	if len(description) > 0 {
		addRequest.Attribute("description", []string{description})
	}
	if len(street) > 0 {
		addRequest.Attribute("street", []string{street})
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

	return s.getOrganizationUnit(conn, &AdEntryFilter{DNs: []string{dn}})
}
