package weather

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

var (
	ErrCityNotExists = errors.New("city not exists")
	ErrCitySameName  = errors.New("there are cities with the same name ")
)

type HeFeng struct {
	ApiKey  string
	BaseUrl string
}

type GeoInfo struct {
	Code     string     `json:"code"`
	Location []Location `json:"location"`
}

type Location struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Lat       string `json:"lat"`
	Lon       string `json:"lon"`
	Adm2      string `json:"adm2"`
	Adm1      string `json:"adm1"`
	Country   string `json:"country"`
	Tz        string `json:"tz"`
	UTCOffset string `json:"utcOffset"`
	IsDst     string `json:"isDst"`
	Type      string `json:"type"`
	Rank      string `json:"rank"`
	FxLink    string `json:"fxLink"`
}

type TopLevel struct {
	Code       string  `json:"code"`
	UpdateTime string  `json:"updateTime"`
	FxLink     string  `json:"fxLink"`
	Daily      []Daily `json:"daily"`
}

type Daily struct {
	FxDate         string `json:"fxDate"`
	Sunrise        string `json:"sunrise"`
	Sunset         string `json:"sunset"`
	Moonrise       string `json:"moonrise"`
	Moonset        string `json:"moonset"`
	MoonPhase      string `json:"moonPhase"`
	MoonPhaseIcon  string `json:"moonPhaseIcon"`
	TempMax        string `json:"tempMax"`
	TempMin        string `json:"tempMin"`
	IconDay        string `json:"iconDay"`
	TextDay        string `json:"textDay"`
	IconNight      string `json:"iconNight"`
	TextNight      string `json:"textNight"`
	Wind360Day     string `json:"wind360Day"`
	WindDirDay     string `json:"windDirDay"`
	WindScaleDay   string `json:"windScaleDay"`
	WindSpeedDay   string `json:"windSpeedDay"`
	Wind360Night   string `json:"wind360Night"`
	WindDirNight   string `json:"windDirNight"`
	WindScaleNight string `json:"windScaleNight"`
	WindSpeedNight string `json:"windSpeedNight"`
	Humidity       string `json:"humidity"`
	Precip         string `json:"precip"`
	Pressure       string `json:"pressure"`
	Vis            string `json:"vis"`
	Cloud          string `json:"cloud"`
	UvIndex        string `json:"uvIndex"`
}

func (h *HeFeng) SearchGeo(city string) (location *Location, err error) {
	client := resty.New()
	geoInfo := &GeoInfo{}
	resp, err := client.R().
		SetPathParam("location", city).
		SetPathParam("key", h.ApiKey).
		SetResult(geoInfo).
		Get("https://geoapi.qweather.com/v2/city/lookup")
	if err != nil {
		//app.Log.Sugar().Error("get city geo info error", err)
		return
	}
	if resp.IsError() {
		//app.Log.Sugar().Error("get city geo info response err")
		return
	}
	if len(geoInfo.Location) == 0 {
		// 城市信息有误
		return nil, ErrCityNotExists
	}
	if len(geoInfo.Location) > 1 {
		// 确认是否当地
		return &geoInfo.Location[0], ErrCitySameName
	}
	return &geoInfo.Location[0], nil
}

func (h *HeFeng) GetCityWeather(city string) []string {
	var result TopLevel
	location, err := h.SearchGeo(city)
	if err != nil && !errors.Is(err, ErrCitySameName) {
		return nil
	}
	if errors.Is(err, ErrCitySameName) {
		//存在同名城市，需要确认
		return nil
	}
	client := resty.New()
	resp, err := client.R().
		SetPathParam("location", location.ID).
		SetPathParam("adm2", location.Adm2).
		SetPathParam("adm1", location.Adm1).
		SetPathParam("key", h.ApiKey).
		SetResult(&result).
		Get("https://devapi.qweather.com/v7/weather/3d")
	if err != nil {
		return nil
	}
	if resp.IsError() {
		//app.Log.Sugar().Error("get city geo info response err")
	}
	weatherDesc := "日期:%s，最高气温:%s，最低气温:%s，白天天气状况:%s，晚间天气状况:%s \n"
	weathers := make([]string, 0)
	for _, v := range result.Daily {
		weathers = append(weathers, fmt.Sprintf(weatherDesc, v.FxDate, v.TempMax, v.TempMin, v.TextDay, v.TextNight))
	}
	return weathers
}
