package job

var locationMap = make(map[string]interface{})

func SetLocationMap(m map[string]interface{}) {
	locationMap = m
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
