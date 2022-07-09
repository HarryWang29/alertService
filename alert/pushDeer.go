package alert

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"time"
)

const pushDeerUrl = "https://api2.pushdeer.com/message/push"

type PushDeer struct {
	Keys []string
}

func (p *PushDeer) Push(title, body string) {
	for _, key := range p.Keys {
		s, err := httplib.Get(pushDeerUrl).
			Param("pushkey", key).
			Param("text", title).
			Param("desp", body).
			Param("type", "markdown").
			SetTimeout(10*time.Second, 10*time.Second).String()
		if err != nil {
			continue
		}
		m := make(map[string]interface{})
		err = json.Unmarshal([]byte(s), &m)
		if err != nil {
			logrus.Error("推送失败")
		}
		var success, total int64
		gjson.Get(s, "content.result").ForEach(func(key, value gjson.Result) bool {
			j := gjson.Parse(value.String())
			count := j.Get("counts").Int()
			total += count
			if j.Get("success").String() == "ok" {
				success += count
			}
			return true
		})
		logrus.Infof("推送成功：%d/%d", success, total)
	}
}
