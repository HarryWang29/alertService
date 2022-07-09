package job

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

var locationMap = make(map[string]interface{})

func init() {
	bs, err := ioutil.ReadFile("ProvinceCode.json")
	if err != nil {
		logrus.Fatalf("ioutil.ReadFile ProvinceCode.json err: %v", err)
	}
	err = json.Unmarshal(bs, &locationMap)
	if err != nil {
		logrus.Fatalf("json.Unmarshal err: %v", err)
	}
}

type Location struct {
	Code     string `json:"code"`
	Province string `json:"province" yaml:"province"`
	City     string `json:"city" yaml:"city"`
	County   string `json:"county" yaml:"county"`
}

func NewLocation(province string, city string, county string) *Location {
	l := &Location{Province: province, City: city, County: county}
	var p, c, cc map[string]interface{}
	if province != "" {
		p = locationMap[province].(map[string]interface{})
		l.Code = p["code"].(string)
	}
	if city != "" {
		c = p[city].(map[string]interface{})
		l.Code = c["code"].(string)
	}
	if county != "" {
		cc = c[county].(map[string]interface{})
		l.Code = cc["code"].(string)
	}
	return l
}
