package config

import (
	"alertService/alert"
	"alertService/job"
	"alertService/store"
	"github.com/pkg/errors"
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

func Load(path string) (*RawConfig, error) {
	c := &RawConfig{}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadFile")
	}

	c, err = ParseConfig(bs)
	if err != nil {
		return nil, errors.Wrap(err, "ParseConfig")
	}

	//初始化日志
	InitLog(c.Log)

	return c, nil
}

func ParseConfig(bs []byte) (*RawConfig, error) {
	c := &RawConfig{}
	err := yaml.Unmarshal(bs, c)
	if err != nil {
		return nil, errors.Wrap(err, "yaml.Unmarshal")
	}
	return c, nil
}

func parseStore(sm map[string]interface{}) (s store.Store, err error) {
	switch sm["type"] {
	case "memory":
		s = store.NewMemory()
		logrus.Info("使用内存缓存")
	case "mysql":
		dns := sm["dns"].(string)
		s = store.NewMysql(dns)
		logrus.Info("使用MySQL缓存")
	default:
		return nil, errors.New("未知的缓存类型")
	}
	return s, nil
}

func parseJobs(jobs []map[string]interface{}, s store.Store, a alert.Alert) (js []job.Job, err error) {
	js = make([]job.Job, 0)
	for _, j := range jobs {
		var jj job.Job
		switch j["type"].(string) {
		case "WeatherWarning":
			jj = job.NewWeatherWarning(j, s, a)
			logrus.Info("开启天气预警")
		default:
			return nil, errors.New("未知的任务类型")
		}
		js = append(js, jj)
	}
	return js, nil
}

func parseAlert(am []map[string]interface{}) (a alert.Alert) {
	if len(am) == 0 {
		return nil
	}
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
