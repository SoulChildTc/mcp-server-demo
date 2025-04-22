package model

type WeatherResponse struct {
	Code string     `json:"code"`
	Now  WeatherNow `json:"now"`
}

type WeatherNow struct {
	ObsTime   string `json:"obsTime"`   // 数据观测时间
	Temp      string `json:"temp"`      // 温度
	FeelsLike string `json:"feelsLike"` // 体感温度
	Icon      string `json:"icon"`      // 天气状况图标代码
	Text      string `json:"text"`      // 天气状况文字描述
	Wind360   string `json:"wind360"`   // 风向360角度
	WindDir   string `json:"windDir"`   // 风向
	WindScale string `json:"windScale"` // 风力等级
	WindSpeed string `json:"windSpeed"` // 风速，公里/小时
	Humidity  string `json:"humidity"`  // 相对湿度
	Precip    string `json:"precip"`    // 当前小时降水量，毫米
	Pressure  string `json:"pressure"`  // 大气压强，百帕
	Vis       string `json:"vis"`       // 能见度，公里
	Cloud     string `json:"cloud"`     // 云量 (可能为空)
	Dew       string `json:"dew"`       // 露点温度 (可能为空)
}

type GeoResponse struct {
	Code     string         `json:"code"`
	Location []LocationInfo `json:"location"`
}

type LocationInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
	Adm1 string `json:"adm1"` // 省份
	Adm2 string `json:"adm2"` // 城市
}
