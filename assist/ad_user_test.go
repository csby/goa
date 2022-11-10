package assist

import (
	"testing"
)

const (
	AdHost     = "192.168.123.101"
	AdBase     = "dc=csby,dc=fun"
	AdPort     = 636
	AdAccount  = "CN=Administrator,CN=Users,DC=csby,DC=fun"
	AdPassword = "Vico0808"
)

func TestAd_Login(t *testing.T) {
	ad := &Ad{
		Host:     AdHost,
		Port:     AdPort,
		Base:     AdBase,
		Account:  AdAccount,
		Password: AdPassword,
	}

	entry, err := ad.Login("dev", "#Dv0808")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("entry: %#v", entry)
	}
}

func TestAd_GetAllUsers(t *testing.T) {
	ad := &Ad{
		Host:     AdHost,
		Port:     AdPort,
		Base:     AdBase,
		Account:  AdAccount,
		Password: AdPassword,
	}
	items, err := ad.GetAllUsers()
	if err != nil {
		t.Error(err)
		return
	}
	c := len(items)
	t.Log("count: ", c)
	i := 0
	for _, item := range items {
		t.Logf("item %03d: %#v", i+1, item)
		i++
	}
}

func TestAd_GetUsers(t *testing.T) {
	ad := &Ad{
		Host:     AdHost,
		Port:     AdPort,
		Base:     AdBase,
		Account:  AdAccount,
		Password: AdPassword,
	}
	parentDN := "OU=组织2,OU=用户账号,DC=csby,DC=fun"

	items, err := ad.GetUsers(parentDN)
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

func TestAd_GetVpnUsers(t *testing.T) {
	ad := &Ad{
		Host:     AdHost,
		Port:     AdPort,
		Base:     AdBase,
		Account:  AdAccount,
		Password: AdPassword,
	}

	items, err := ad.GetVpnUsers()
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
