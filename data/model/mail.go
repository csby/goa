package model

type MailAddress struct {
	Address string `json:"address" required:"true" note:"地址"`
}

type MailCommitter struct {
	MailAddress

	CommitBy string `json:"commitBy" note:"提交者"`
}

type MailBlockAddress struct {
	MailCommitter

	BlockCount         int64  `json:"blockCount" note:"拒收次数"`
	CommitDateTime     string `json:"commitDateTime" note:"提交时间"'`
	FirstBlockDateTime string `json:"firstBlockDateTime" note:"首次阻止时间"`
	LastBlockDateTime  string `json:"lastBlockDateTime" note:"最后阻止时间"`
}

type MailBlockIP struct {
	Ip                  string `json:"ip" note:"IP地址"`
	RejectCount         int64  `json:"rejectCount" note:"阻止次数"`
	FirstRejectDateTime string `json:"firstRejectDateTime" note:"首次阻止时间"`
	LastRejectDateTime  string `json:"lastRejectDateTime" note:"最后阻止时间"`
	IpLocator           string `json:"ipLocator" note:"IP地址归属地"`
}
