package model

import "strings"

type DhcpFilter struct {
	Address string `json:"address" required:"true" note:"MAC地址"`
	Owner   string `json:"owner" note:"所有者"`
	Type    string `json:"type" note:"设备类型"`
	Remark  string `json:"remark" note:"备注信息"`
}

type DhcpFilterCollection []*DhcpFilter

func (s DhcpFilterCollection) Len() int {
	return len(s)
}

func (s DhcpFilterCollection) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s DhcpFilterCollection) Less(i, j int) bool {
	owner := strings.Compare(s[i].Owner, s[j].Owner)
	if owner == 0 {
		return strings.Compare(s[i].Type, s[j].Type) < 0
	} else {
		return owner < 0
	}
}

type DhcpFilterAll struct {
	items DhcpFilterCollection
}

func (s *DhcpFilterAll) Add(address, comment string) *DhcpFilter {
	if len(address) < 1 {
		return nil
	}

	filter := &DhcpFilter{
		Address: address,
	}

	fields := make([]string, 0)
	comments := strings.Split(comment, "-")
	c := len(comments)
	for i := 0; i < c; i++ {
		item := strings.TrimSpace(comments[i])
		fields = append(fields, item)
	}
	c = len(fields)
	if c > 0 {
		filter.Owner = fields[0]
	}
	if c > 1 {
		filter.Type = fields[1]
	}
	if c > 2 {
		filter.Remark = fields[2]
	}

	if s.items == nil {
		s.items = make(DhcpFilterCollection, 0)
	}
	s.items = append(s.items, filter)

	return filter
}

func (s *DhcpFilterAll) Items(owner string) DhcpFilterCollection {
	items := make(DhcpFilterCollection, 0)

	if len(owner) < 1 {
		items = s.items
	} else {
		c := len(s.items)
		for i := 0; i < c; i++ {
			item := s.items[i]
			if item == nil {
				continue
			}
			if strings.Contains(item.Owner, owner) {
				items = append(items, item)
			}
		}
	}

	return items
}

type DhcpFilterListArgument struct {
	Owner string `json:"owner" note:"所有者"`
}

type DhcpFilterDeleteArgument struct {
	Address string `json:"address" required:"true" note:"MAC地址"`
}

type DhcpFilterModifyArgument struct {
	Address string `json:"address" required:"true" note:"原MAC地址"`

	Filter DhcpFilter `json:"filter" note:"新筛选器"`
}

type DhcpLease struct {
	IpV4    string `json:"ipV4" note:"IPv4地址"`
	Address string `json:"address" note:"MAC地址"`
	Comment string `json:"comment" note:"描述"`
}
