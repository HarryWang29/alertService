package config

import (
	"github.com/HarryWang29/alertService/job"
	"github.com/pkg/errors"
)

type App struct {
	Jobs []job.Job
}

func NewApp() *App {
	return &App{}
}

func (a *App) SetJobs(jobs []job.Job) {
	a.Jobs = jobs
}

func (a *App) Parse(config *RawConfig) (err error) {
	//初始化存储
	s, err := parseStore(config.Store)
	if err != nil {
		return errors.Wrap(err, "parseStore")
	}
	//初始化推送服务
	alert := parseAlert(config.Alert)
	//初始化job
	a.Jobs, err = parseJobs(config.Jobs, s, alert)
	if err != nil {
		return errors.Wrap(err, "parseJobs")
	}
	return nil
}
