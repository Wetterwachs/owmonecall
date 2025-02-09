package owmonecall

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type RainData struct {
	OneHour float64 `json:"1h"`
}

type SnowData struct {
	OneHour float64 `json:"1h"`
}

type ShortWeatherData struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type CurrentWeatherData struct {
	Dt          int64              `json:"dt"`
	SunRise     int64              `json:"sunrise"`
	SunSet      int64              `json:"sunset"`
	Temperature float64            `json:"temp"`
	FeelsLike   float64            `json:"feels_like"`
	Pressure    int                `json:"pressure"`
	Humidity    int                `json:"humidity"`
	DewPoint    float64            `json:"dew_point"`
	Clouds      int                `json:"clouds"`
	UVIndex     int                `json:"uvi"`
	Visibility  int                `json:"visibility"`
	WindSpeed   float64            `json:"wind_speed"`
	WindGust    float64            `json:"wind_gust"`
	WindDegree  int                `json:"wind_deg"`
	Rain        RainData           `json:"rain"`
	Snow        SnowData           `json:"snow"`
	Weather     []ShortWeatherData `json:"weather"`
}

type MinutelyWeatherData struct {
	Dt            int64 `json:"dt"`
	Precipitation int   `json:"precipitation"`
}

type HourlyWeatherData struct {
	Dt          int64              `json:"dt"`
	Temperature float64            `json:"temp"`
	FeelsLike   float64            `json:"feels_like"`
	Pressure    int                `json:"pressure"`
	Humidity    int                `json:"humidity"`
	DewPoint    float64            `json:"dew_point"`
	UVIndex     float64            `json:"uvi"`
	Clouds      int                `json:"clouds"`
	Visibility  int                `json:"visibility"`
	WindSpeed   float64            `json:"wind_speed"`
	WindGust    float64            `json:"wind_gust"`
	WindDegree  int                `json:"wind_deg"`
	Propability float64            `json:"pop"`
	Rain        RainData           `json:"rain"`
	Snow        SnowData           `json:"snow"`
	Weather     []ShortWeatherData `json:"weather"`
}

type DailyTemperatureData struct {
	Morning float64 `json:"morn"`
	Day     float64 `json:"day"`
	Eve     float64 `json:"eve"`
	Night   float64 `json:"night"`
	Minimum float64 `json:"min"`
	Maximum float64 `json:"max"`
}

type DailyFeelsLikeData struct {
	Morning float64 `json:"morn"`
	Day     float64 `json:"day"`
	Eve     float64 `json:"eve"`
	Night   float64 `json:"night"`
}

type DailyWeatherData struct {
	Dt          int64                `json:"dt"`
	SunRise     int64                `json:"sunrise"`
	SunSet      int64                `json:"sunset"`
	MoonRise    int64                `json:"moonrise"`
	MooSet      int64                `json:"moonset"`
	MoonPhase   float64              `json:"moon_phase"`
	Temperature DailyTemperatureData `json:"temp"`
	FeelsLike   DailyFeelsLikeData   `json:"feels_like"`
	Pressure    int                  `json:"pressure"`
	Humidity    int                  `json:"humidity"`
	DewPoint    float64              `json:"dew_point"`
	WindSpeed   float64              `json:"wind_speed"`
	WindGust    float64              `json:"wind_gust"`
	WindDegree  int                  `json:"wind_deg"`
	Clouds      int                  `json:"clouds"`
	UVIndex     float64              `json:"uvi"`
	Propability float64              `json:"pop"`
	Rain        float64              `json:"rain"`
	Snow        float64              `json:"snow"`
	Weather     []ShortWeatherData   `json:"weather"`
}

type AlertData struct {
	SenderName  string `json:"sender_name"`
	Event       string `json:"event"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
}

type ConfigDataExcludes struct {
	Current  string
	Minutely string
	Hourly   string
	Daily    string
	Alerts   string
}

type ConfigData struct {
	Latitude  float64
	Longitude float64
	ApiKey    string
	Excludes  *ConfigDataExcludes
	Units     string
	Language  string
}

type WeatherData struct {
	Latitude       float64               `json:"lat"`
	Longitude      float64               `json:"lon"`
	Timezone       string                `json:"timezone"`
	TimezoneOffset int64                 `json:"timezone_offset"`
	Current        CurrentWeatherData    `json:"current"`
	Minutely       []MinutelyWeatherData `json:"minutely"`
	Hourly         []HourlyWeatherData   `json:"hourly"`
	Daily          []DailyWeatherData    `json:"daily"`
	Alerts         []AlertData           `json:"alerts"`
	config         *ConfigData
	client         *http.Client
}

const (
	Current  = "current"
	Minutely = "minutely"
	Hourly   = "hourly"
	Daily    = "daily"
	Alerts   = "alerts"
	Standard = "standard"
	Metric   = "metric"
	Imperial = "imperial"
)

func NewWeatherData(latitude float64, longitude float64, dataToGet string, units string, language string, apiKey string) (w *WeatherData, err error) {
	w = &WeatherData{}
	w.client = http.DefaultClient
	return w, w.Configure(latitude, longitude, dataToGet, units, language, apiKey)
}

func (w *WeatherData) Configure(latitude float64, longitude float64, dataToGet string, units string, language string, apiKey string) (err error) {
	w.config, err = NewConfigData(latitude, longitude, dataToGet, units, language, apiKey)
	return err
}

func (w *WeatherData) Update() (err error) {
	response, err := w.client.Get(w.config.GetURL())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		if err := json.NewDecoder(response.Body).Decode(&w); err != nil {
			return err
		}

	}

	return nil
}

func NewConfigData(latitude float64, longitude float64, dataToGet string, units string, language string, apiKey string) (c *ConfigData, err error) {
	c = &ConfigData{}
	c.Configure(latitude, longitude, dataToGet, units, language, apiKey)
	return c, nil
}

func (c *ConfigData) Configure(latitude float64, longitude float64, dataToGet string, units string, language string, apiKey string) (err error) {
	c.Latitude = latitude
	c.Longitude = longitude
	c.Excludes, err = NewConfigDataExcludes(dataToGet)
	c.Units = units
	c.Language = language
	c.ApiKey = apiKey
	return err
}

func (c *ConfigData) GetURL() (url string) {
	url = "https://api.openweathermap.org/data/3.0/onecall?"
	url += "lat=" + strconv.FormatFloat(c.Latitude, 'f', 2, 64)
	url += "&lon=" + strconv.FormatFloat(c.Longitude, 'f', 2, 64)
	excludes := c.Excludes.GetExclude()
	if excludes != "" {
		url += "&exclude=" + excludes
	}
	if c.Units != "" {
		url += "&units=" + c.Units
	}
	if c.Language != "" {
		url += "&lang=" + c.Language
	}
	url += "&appid=" + c.ApiKey

	return url
}

func NewConfigDataExcludes(settings string) (c *ConfigDataExcludes, err error) {
	c = &ConfigDataExcludes{}
	return c, c.Configure(settings)
}

func (c *ConfigDataExcludes) Configure(settings string) (err error) {
	if !strings.Contains(settings, Current) {
		c.Current = Current
	}
	if !strings.Contains(settings, Minutely) {
		c.Minutely = Minutely
	}
	if !strings.Contains(settings, Hourly) {
		c.Hourly = Hourly
	}
	if !strings.Contains(settings, Daily) {
		c.Daily = Daily
	}
	if !strings.Contains(settings, Alerts) {
		c.Alerts = Alerts
	}

	return nil
}

func (c *ConfigDataExcludes) GetExclude() (excludes string) {
	var txt string
	if c.Current != "" {
		txt = c.Current + ","
	}
	if c.Minutely != "" {
		txt += c.Minutely + ","
	}
	if c.Hourly != "" {
		txt += c.Hourly + ","
	}
	if c.Daily != "" {
		txt += c.Daily + ","
	}
	if c.Alerts != "" {
		txt += c.Alerts
	}
	return txt
}
