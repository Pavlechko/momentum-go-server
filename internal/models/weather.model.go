package models

type WeatherInput struct {
	City   string `json:"city" binding:"required"`
	Source string `json:"source" binding:"required"`
}

type WeatherFild struct {
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type OpenWeatherResponse struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed     float64 `json:"speed"`
		Direction float64 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Sys struct {
		Country string `json:"Country"`
	} `json:"sys"`
	Weather []WeatherFild `json:"weather"`
	Name    string        `json:"name"`
}

type TomorrowWeatherResponse struct {
	Data struct {
		Values struct {
			Temperature         float64 `json:"temperature"`
			TemperatureApparent float64 `json:"temperatureApparent"`
			Humidity            int     `json:"humidity"`
			WindSpeed           float64 `json:"windSpeed"`
			WindDirection       float64 `json:"windDirection"`
			CloudCover          int     `json:"cloudCover"`
			RainIntensity       float64 `json:"rainIntensity"`
			SnowIntensity       float64 `json:"snowIntensity"`
		} `json:"values"`
	} `json:"data"`
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
}

type FrontendWeatherResponse struct {
	Temp       float64 `json:"temp"`
	FeelsLike  float64 `json:"feels_like"`
	Humidity   int     `json:"humidity"`
	WindSpeed  string  `json:"wind_speed"`
	WeaterMain string  `json:"main"`
	WeaterIcon string  `json:"icon"`
	City       string  `json:"city"`
	Source     string  `json:"source"`
}

var CITIES = []string{
	"Kyiv", "Kharkiv", "Odesa", "Dnipro", "Donetsk", "Lviv", "Vinnytsia", "Zaporizhia", "Ivano-Frankivsk",
	"Poltava", "Mykolaiv", "Chernihiv", "Cherkasy", "Chernivtsi", "Sumy", "Kherson", "Luhansk", "Rivne",
	"Ternopil", "Zhytomyr", "Khmelnytskyi", "Uzhhorod", "Lutsk", "Simferopol", "Kramatorsk", "Kropyvnytskyi",
}

var WEATHER_PROVIDERS = []string{"OpenWeather", "TomorrowWeather"}
