package main

import (
	"fmt"
	"github.com/csby/goa/config"
	"github.com/csby/gwsf/gcfg"
	"github.com/csby/gwsf/glog"
	"github.com/csby/gwsf/gserver"
	"github.com/csby/gwsf/gtype"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	moduleType    = "server"
	moduleName    = "goa"
	moduleRemark  = "办公自动化系统"
	moduleVersion = "1.0.1.0"
)

var (
	cfg              = config.NewConfig()
	log              = &glog.Writer{Level: glog.LevelAll}
	svr gtype.Server = nil
)

func init() {
	moduleArgs := &gtype.Args{}
	serverArgs := &gtype.SvcArgs{}
	moduleArgs.Parse(os.Args, moduleType, moduleName, moduleVersion, moduleRemark, serverArgs)
	now := time.Now()
	cfg.Module.Type = moduleType
	cfg.Module.Name = moduleName
	cfg.Module.Version = moduleVersion
	cfg.Module.Remark = moduleRemark
	cfg.Module.Path = moduleArgs.ModulePath()
	cfg.Svc.BootTime = now

	rootFolder := filepath.Dir(moduleArgs.ModuleFolder())
	cfgFolder := filepath.Join(rootFolder, "cfg")
	cfgName := fmt.Sprintf("%s.json", moduleName)
	if serverArgs.Help {
		serverArgs.ShowHelp(cfgFolder, cfgName)
		os.Exit(11)
	}

	if serverArgs.Pkg {
		pkg := &Pkg{binPath: cfg.Module.Path}
		pkg.Run()
		os.Exit(0)
	}

	// init config
	svcArgument := ""
	cfgPath := serverArgs.Cfg
	if cfgPath != "" {
		svcArgument = fmt.Sprintf("-cfg=%s", cfgPath)
	} else {
		cfgPath = filepath.Join(cfgFolder, cfgName)
	}
	_, err := os.Stat(cfgPath)
	if os.IsNotExist(err) {
		err = cfg.SaveToFile(cfgPath)
		if err != nil {
			fmt.Println("generate configure file fail: ", err)
		}
	} else {
		err = cfg.LoadFromFile(cfgPath)
		if err != nil {
			fmt.Println("load configure file fail: ", err)
		}
	}
	cfg.Path = cfgPath
	cfg.Load = cfg.DoLoad
	cfg.Save = cfg.DoSave
	cfg.InitId()

	//cfg.SaveToFile(cfgPath)

	// init certificate
	if cfg.Https.Enabled {
		certFilePath := cfg.Https.Cert.Server.File
		if certFilePath == "" {
			certFilePath = filepath.Join(rootFolder, "crt", "server.pfx")
			cfg.Https.Cert.Server.File = certFilePath
		}
	}

	// init path of site
	if cfg.Site.Root.Path == "" {
		cfg.Site.Root.Path = filepath.Join(rootFolder, "site", "root")
	}
	if cfg.Site.Doc.Path == "" {
		cfg.Site.Doc.Path = filepath.Join(rootFolder, "site", "doc")
	}
	if cfg.Site.Opt.Path == "" {
		cfg.Site.Opt.Path = filepath.Join(rootFolder, "site", "opt")
	}
	oaSiteExisted := false
	for i := 0; i < len(cfg.Site.Apps); i++ {
		site := cfg.Site.Apps[i]
		if site.Uri == config.StaffSiteUri {
			oaSiteExisted = true
			break
		}
	}
	if !oaSiteExisted {
		if cfg.Site.Apps == nil {
			cfg.Site.Apps = make([]gcfg.SiteApp, 0)
		}
		cfg.Site.Apps = append(cfg.Site.Apps, gcfg.SiteApp{
			Name: config.StaffSiteName,
			Uri:  config.StaffSiteUri,
		})
	}

	// init uri for dhcp filter api uri
	if cfg.Dhcp.Api.Uri.Filter.List == "" {
		cfg.Dhcp.Api.Uri.Filter.List = "/api/dhcp/filter/list"
	}
	if cfg.Dhcp.Api.Uri.Filter.Add == "" {
		cfg.Dhcp.Api.Uri.Filter.Add = "/api/dhcp/filter/add"
	}
	if cfg.Dhcp.Api.Uri.Filter.Del == "" {
		cfg.Dhcp.Api.Uri.Filter.Del = "/api/dhcp/filter/del"
	}
	if cfg.Dhcp.Api.Uri.Filter.Mod == "" {
		cfg.Dhcp.Api.Uri.Filter.Mod = "/api/dhcp/filter/mod"
	}
	if cfg.Dhcp.Api.Uri.Lease.List == "" {
		cfg.Dhcp.Api.Uri.Lease.List = "/api/dhcp/lease/list"
	}

	// init uri for svn api uri
	if cfg.Svn.Api.Uri.Repository.Add == "" {
		cfg.Svn.Api.Uri.Repository.Add = "/api/svn/repository/new"
	}
	if cfg.Svn.Api.Uri.Repository.List == "" {
		cfg.Svn.Api.Uri.Repository.List = "/api/svn/repository/list"
	}
	if cfg.Svn.Api.Uri.Repository.Folder == "" {
		cfg.Svn.Api.Uri.Repository.Folder = "/api/svn/folder/list"
	}
	if cfg.Svn.Api.Uri.Permission.List == "" {
		cfg.Svn.Api.Uri.Permission.List = "/api/svn/permission/list"
	}
	if cfg.Svn.Api.Uri.Permission.Add == "" {
		cfg.Svn.Api.Uri.Permission.Add = "/api/svn/permission/add"
	}
	if cfg.Svn.Api.Uri.Permission.Mod == "" {
		cfg.Svn.Api.Uri.Permission.Mod = "/api/svn/permission/mod"
	}
	if cfg.Svn.Api.Uri.Permission.Del == "" {
		cfg.Svn.Api.Uri.Permission.Del = "/api/svn/permission/del"
	}
	if cfg.Svn.Api.Uri.User.Permission == "" {
		cfg.Svn.Api.Uri.User.Permission = "/api/svn/user/permission/list"
	}

	// init path of system service
	if cfg.Sys.Svc.Custom.App == "" {
		cfg.Sys.Svc.Custom.App = filepath.Join(rootFolder, "svc", "custom")
	}
	if cfg.Sys.Svc.Custom.Log == "" {
		cfg.Sys.Svc.Custom.Log = filepath.Join(rootFolder, "log", "svc", "custom")
	}

	// init service
	if strings.TrimSpace(cfg.Svc.Name) == "" {
		cfg.Svc.Name = moduleName
	}
	cfg.Svc.Args = svcArgument
	svcName := cfg.Svc.Name
	log.Init(cfg.Log.Level, svcName, cfg.Log.Folder)
	hdl := NewHandler(log, cfg)
	svr, err = gserver.NewServer(log, &cfg.Config, hdl)
	if err != nil {
		fmt.Println("init service fail: ", err)
		os.Exit(12)
	}
	if !svr.Interactive() {
		cfg.Svc.Restart = svr.Restart
	}
	serverArgs.Execute(svr)

	// information
	log.Std = true
	zoneName, zoneOffset := now.Zone()
	LogInfo("start at: ", moduleArgs.ModulePath())
	LogInfo("run as service: ", !svr.Interactive())
	LogInfo("version: ", moduleVersion)
	LogInfo("zone: ", zoneName, "-", zoneOffset/int(time.Hour.Seconds()))
	LogInfo("log path: ", cfg.Log.Folder)
	LogInfo("log level: ", cfg.Log.Level)
	LogInfo("configure path: ", cfgPath)
	LogInfo("configure info: ", cfg)
}
