package assist

import (
	"testing"
)

func TestAd_GetOrganizationUnits(t *testing.T) {
	ad := &Ad{
		Host:     AdHost,
		Port:     AdPort,
		Base:     AdBase,
		Account:  AdAccount,
		Password: AdPassword,
	}
	parentDN := "OU=用户账号,DC=csby,DC=fun"

	items, err := ad.GetOrganizationUnits(parentDN)
	if err != nil {
		t.Error(err)
		return
	}
	c := len(items)
	t.Log("count: ", c)
	for i := 0; i < c; i++ {
		item := items[i]
		t.Logf("item %03d: %#v", i+1, item)
	}
}

func TestAd_AddOrganizationUnit(t *testing.T) {
	ad := &Ad{
		Host:     AdHost,
		Port:     AdPort,
		Base:     AdBase,
		Account:  AdAccount,
		Password: AdPassword,
	}
	dn := "OU=TestOU,DC=csby,DC=fun"

	item, err := ad.AddOrganizationUnit(dn, "描述", "街道")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("item: %#v", item)
}
