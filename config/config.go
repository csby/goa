package config

import (
	"encoding/json"
	"fmt"
	"github.com/csby/gwsf/gcfg"
	"github.com/csby/gwsf/gtype"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	sync.RWMutex
	gcfg.Config

	Auth Auth `json:"auth" note:"授权服务"`
	Dhcp Dhcp `json:"dhcp" note:"DHCP服务"`
	Svn  Svn  `json:"svn" note:"SVN服务"`
	Ad   MsAd `json:"ad" note:"AD服务"`
	Mail Mail `json:"mail" note:"邮件服务"`
}

func NewConfig() *Config {
	return &Config{
		Config: gcfg.Config{
			Log: gcfg.Log{
				Folder: "",
				Level:  "error|warning|info",
			},
			Http: gcfg.Http{
				Enabled:     true,
				Port:        8085,
				BehindProxy: false,
			},
			Https: gcfg.Https{
				Enabled:     false,
				Port:        8443,
				BehindProxy: false,
				Cert: gcfg.Crt{
					Ca: gcfg.CrtCa{
						File: "",
					},
					Server: gcfg.CrtPfx{
						File:     "",
						Password: "",
					},
				},
				RequestClientCert: false,
			},
			Site: gcfg.Site{
				Doc: gcfg.SiteDoc{
					Enabled: true,
				},
				Opt: gcfg.SiteOpt{
					Users: []*gcfg.SiteOptUser{
						{
							Account:  "admin",
							Password: "admin",
							Name:     "内置管理员",
						},
					},
				},
				Apps: []gcfg.SiteApp{
					{
						Name: StaffSiteName,
						Uri:  StaffSiteUri,
					},
				},
			},
			ReverseProxy: gcfg.Proxy{
				Enabled: false,
				Servers: []*gcfg.ProxyServer{
					{
						Id:      gtype.NewGuid(),
						Name:    "http",
						Disable: true,
						TLS:     false,
						IP:      "",
						Port:    "80",
						Targets: []*gcfg.ProxyTarget{},
					},
					{
						Id:      gtype.NewGuid(),
						Name:    "https",
						Disable: true,
						TLS:     true,
						IP:      "",
						Port:    "443",
						Targets: []*gcfg.ProxyTarget{},
					},
				},
			},
			Sys: gcfg.System{
				Svc: gcfg.Service{
					Enabled: false,
					Tomcats: []*gcfg.ServiceTomcat{},
					Others:  []*gcfg.ServiceOther{},
					Nginxes: []*gcfg.ServiceNginx{},
					Files:   []*gcfg.ServiceFile{},
				},
			},
		},
		Dhcp: Dhcp{
			Api: DhcpApi{
				Url: "http://192.168.123.101:8085",
				Uri: DhcpApiUri{
					Filter: DhcpApiUriFilter{},
					Lease:  DhcpApiUriLease{},
				},
			},
		},
		Svn: Svn{
			Api: SvnApi{
				Url: "http://192.168.123.101:8085",
			},
		},
		Ad: MsAd{
			Host: "127.0.0.1",
			Port: 636,
			Base: "DC=example,DC=com",
			Account: MsAdAccount{
				Account:  "CN=Administrator,CN=Users,DC=example,DC=com",
				Password: "",
			},
			AdminGroup: "oa.admins",
			Root: MsAdRoot{
				Server: "OU=服务器,DC=example,DC=com",
				Share:  "OU=共享目录,DC=example,DC=com",
				User:   "OU=用户账号,DC=example,DC=com",
				Svn:    "OU=SVN,DC=example,DC=com",
			},
		},
		Mail: Mail{
			Api: MailApi{
				Url: "http://172.16.100.3:11034",
			},
		},
	}
}

func (s *Config) LoadFromFile(filePath string) error {
	s.Lock()
	defer s.Unlock()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

func (s *Config) SaveToFile(filePath string) error {
	s.Lock()
	defer s.Unlock()

	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	fileFolder := filepath.Dir(filePath)
	_, err = os.Stat(fileFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(fileFolder, 0777)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, string(bytes[:]))

	return err
}

func (s *Config) DoLoad() (*gcfg.Config, error) {
	c := &Config{}
	e := c.LoadFromFile(s.Path)
	if e != nil {
		return nil, e
	}

	return &c.Config, nil
}

func (s *Config) DoSave(cfg *gcfg.Config) error {
	if cfg == nil {
		return nil
	}

	c := &Config{}
	e := c.LoadFromFile(s.Path)
	if e != nil {
		return e
	}

	c.Config = *cfg
	return c.SaveToFile(s.Path)
}

func (s *Config) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

func (s *Config) FormatString() string {
	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes[:])
}
