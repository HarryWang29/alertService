package config

import (
	"alertService/alert"
	"alertService/job"
	"alertService/store"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type RawConfig struct {
	Jobs  []map[string]interface{} `yaml:"jobs"`
	Store map[string]interface{}   `yaml:"store"`
	Alert []map[string]interface{} `yaml:"alert"`
	Log   *LogConfig               `yaml:"log"`
}

func Load(path string) *RawConfig {
	c := &RawConfig{}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatalf("ioutil.ReadFile err: %v", err)
	}
	err = yaml.Unmarshal(bs, c)
	if err != nil {
		logrus.Fatalf("yaml.Unmarshal error: %v", err)
	}
	return c
}

func parseStore(sm map[string]interface{}) store.Store {
	var s store.Store
	switch sm["type"] {
	case "memory":
		s = store.NewMemory()
		logrus.Info("使用内存缓存")
	case "mysql":
		dns := sm["dns"].(string)
		s = store.NewMysql(dns)
		logrus.Info("使用MySQL缓存")
	}
	return s
}

func parseJobs(jobs []map[string]interface{}, s store.Store, a alert.Alert) []job.Job {
	js := make([]job.Job, 0)
	for _, j := range jobs {
		var jj job.Job
		switch j["type"].(string) {
		case "WeatherWarning":
			jj = job.NewWeatherWarning(j, s, a)
			logrus.Info("开启天气预警")
		}
		js = append(js, jj)
	}
	return js
}

func parseAlert(am []map[string]interface{}) alert.Alert {
	var a alert.Alert
	switch am[0]["type"] {
	case "pushDeer":
		pd := &alert.PushDeer{}
		for _, k := range am[0]["keys"].([]interface{}) {
			pd.Keys = append(pd.Keys, k.(string))
		}
		a = pd
		logrus.Info("使用PushDeer推送")
	}
	return a
}
