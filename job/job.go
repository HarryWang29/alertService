package job

import (
	"alertService/alert"
	"alertService/store"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type JobBase struct {
	Name     string `yaml:"name"`
	Schedule string `yaml:"schedule"`
	Store    store.Store
	Alert    alert.Alert
}

func (j *JobBase) Run() {
	println("JobBase")
}

func (j *JobBase) GetName() string {
	return j.Name
}

func (j *JobBase) GetSchedule() string {
	return j.Schedule
}

type Job interface {
	Run()
	GetName() string
	GetSchedule() string
}

type Schedule struct {
	cron  *cron.Cron
	store store.Store
}

func Start(jobs []Job) {
	s := &Schedule{
		cron: cron.New(),
	}
	for _, job := range jobs {
		_, err := s.cron.AddJob(job.GetSchedule(), job)
		if err != nil {
			logrus.Fatalf("s.cron.AddJob err: %v", err)
		}
	}
	s.cron.Run()
}
