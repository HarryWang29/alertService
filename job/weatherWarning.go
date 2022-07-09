package job

import (
	"alertService/alert"
	"alertService/store"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/iris-contrib/schema"
	"github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const weatherWarningUrl = "http://www.nmc.cn/rest/findAlarm"

type WeatherWarning struct {
	*JobBase
	CodeLocations     map[string]*Location
	ProvinceLocations map[string]*Location
	Title             string
}

func NewWeatherWarning(config map[string]interface{}, s store.Store, a alert.Alert) *WeatherWarning {
	//加载配置
	weather := &WeatherWarning{
		JobBase: &JobBase{
			Name:     config["type"].(string),
			Schedule: config["schedule"].(string),
			Store:    s,
			Alert:    a,
		},
		Title:             config["title"].(string),
		CodeLocations:     make(map[string]*Location),
		ProvinceLocations: make(map[string]*Location),
	}
	locations := config["location"].([]interface{})
	if len(locations) == 0 {
		ll := NewLocation("上海市", "", "")
		weather.CodeLocations[ll.Code] = ll
		weather.ProvinceLocations[ll.Province] = ll
		logrus.Info("默认订阅天气预警：上海市")
		return weather
	}
	//查询映射
	for _, location := range locations {
		locationMap := location.(map[interface{}]interface{})
		var province, city, county string
		if locationMap["province"] != nil {
			province = locationMap["province"].(string)
		}
		if locationMap["city"] != nil {
			city = locationMap["city"].(string)
		}
		if locationMap["county"] != nil {
			county = locationMap["county"].(string)
		}

		ll := NewLocation(province, city, county)
		weather.CodeLocations[ll.Code] = ll
		weather.ProvinceLocations[ll.Province] = ll
		logrus.Infof("已订阅天气预警：%s/%s/%s", province, city, county)
	}
	return weather
}

type weatherWarningReq struct {
	PageNo      int    `url:"pageNo"`
	PageSize    int    `url:"pageSize"`
	SignalType  string `url:"signaltype"`
	SignalLevel string `url:"signallevel"`
	Province    string `url:"province"`
	TimeStamp   int64  `url:"_"`
}

type weatherWarningResp struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		Page struct {
			PageNo    int             `json:"pageNo"`
			PageSize  int             `json:"pageSize"`
			Count     int             `json:"count"`
			Prev      int             `json:"prev"`
			Next      int             `json:"next"`
			List      []*weatherAlert `json:"list"`
			TotalPage int             `json:"totalPage"`
		} `json:"page"`
		ProvinceAlarms []interface{} `json:"provinceAlarms"`
		Stat           struct {
			City struct {
				B int `json:"b"`
				O int `json:"o"`
				R int `json:"r"`
				Y int `json:"y"`
			} `json:"city"`
			County struct {
				B int `json:"b"`
				O int `json:"o"`
				R int `json:"r"`
				Y int `json:"y"`
			} `json:"county"`
			Province struct {
				B int `json:"b"`
				O int `json:"o"`
				R int `json:"r"`
				Y int `json:"y"`
			} `json:"province"`
		} `json:"stat"`
	} `json:"data"`
}

type weatherAlert struct {
	Alertid   string `json:"alertid"`
	Issuetime string `json:"issuetime"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	Pic       string `json:"pic"`
}

func (w *WeatherWarning) Run() {
	//获取天气预警
	list := make([]*weatherAlert, 0)
	//控制并发数量
	worker := make(chan struct{}, 10)
	//收集结果
	ret := make(chan *weatherAlert, 1000)
	done := false
	go func(done *bool) {
		for item := range ret {
			list = append(list, item)
		}
		*done = true
	}(&done)
	//等待结束
	wg := sync.WaitGroup{}
	for _, location := range w.ProvinceLocations {
		wg.Add(1)
		worker <- struct{}{}
		go func(location *Location, wg *sync.WaitGroup, ret chan<- *weatherAlert) {
			defer func() {
				wg.Done()
				<-worker
			}()
			list, err := w.getAlert(location)
			if err != nil {
				return
			}
			for _, a := range list {
				code := a.Alertid[0:6]
				alertID := a.Alertid[15:]
				iAlertID, err := strconv.ParseInt(alertID, 10, 64)
				if err != nil {
					logrus.Errorf("strconv.ParseInt(%s) error: %v", alertID, err)
					continue
				}
				//检查是否订阅
				if !w.CheckSub(a) {
					continue
				}
				k, _ := w.Store.GetLast(code)
				var iK int64
				if k != "" {
					iK, err = strconv.ParseInt(k, 10, 64)
					if err != nil {
						logrus.Errorf("strconv.ParseInt(%s) error: %v", k, err)
						continue
					}
				}
				//若此地区已经推送，则跳过
				if iK >= iAlertID {
					continue
				}
				_ = w.Store.Set(code, alertID, a)
				ret <- a
			}
		}(location, &wg, ret)
	}
	wg.Wait()
	close(ret)
	//等待worker结束
	for !done {
		time.Sleep(10 * time.Millisecond)
	}
	if len(list) == 0 {
		return
	}
	//如果有新增，则发送推送
	urlBase, err := url.Parse(weatherWarningUrl)
	if err != nil {
		logrus.Errorf("url.Parse error: %v", err)
		return
	}

	pushBody := ""
	for _, v := range list {
		urlBase.Path = v.Url
		pushBody += fmt.Sprintf("[%s](%s)\t%s\n\n", v.Title, urlBase.String(), v.Issuetime)
	}
	if len(pushBody) == 0 {
		return
	}
	//发送推送
	fmt.Println(pushBody)
	w.Alert.Push(w.Title, pushBody)
}

func (w *WeatherWarning) getAlert(location *Location) (list []*weatherAlert, err error) {
	if location == nil {
		return
	}
	//获取最新的天气预警
	req := &weatherWarningReq{
		PageNo:      1,
		PageSize:    100,
		SignalType:  "",
		SignalLevel: "",
		Province:    location.Province,
		TimeStamp:   time.Now().UnixMilli(),
	}
	v := url.Values{}
	err = schema.NewEncoder().Encode(req, v)
	if err != nil {
		return
	}
	resp := &weatherWarningResp{}
	err = httplib.Get(fmt.Sprintf("%s?%s", weatherWarningUrl, v.Encode())).
		SetTimeout(time.Second*10, time.Second*10).ToJSON(resp)
	if err != nil {
		return
	}

	return resp.Data.Page.List, nil
}

func (w *WeatherWarning) CheckSub(alert *weatherAlert) bool {
	for key := range w.CodeLocations {
		var i int
		for i = len(key) - 1; i >= 0; i-- {
			if key[i] != '0' {
				break
			}
		}
		key = key[:i+1]
		if strings.HasPrefix(alert.Alertid, key) {
			return true
		}
	}
	return false
}
