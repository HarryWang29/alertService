package config

import (
	"alertService/job"
	"os"
)

type App struct {
	Jobs []job.Job
}

func New() *App {
	//加载配置文件
	configPath := ""
	if os.Getenv("alert_service") == "dev" {
		configPath = "./config_dev.yaml"
	} else {
		configPath = "./config.yaml"
	}

	//加载配置文件
	config := Load(configPath)
	//初始化日志
	InitLog(config.Log)

	//初始化存储
	s := parseStore(config.Store)
	//初始化推送服务
	a := parseAlert(config.Alert)
	//初始化job
	app := &App{}
	app.Jobs = parseJobs(config.Jobs, s, a)
	return app
}
