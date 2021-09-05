package main

import (
	"fmt"
	"os"
	owm "github.com/Wetterwachs/owmonecall"
	"time"
)

const APIKEY = ""

func main() {
	w, err := owm.NewWeatherData(
		50.95,                              // Latitude
		6.95,                              // Longitude
		owm.Current+owm.Hourly+owm.Daily,   // data to get
		owm.Metric,                         // units
		"DE",                               // Language
		APIKEY) // ApiKey

	if err != nil {
		fmt.Println("Error while confiure (", err, ")")
		os.Exit(-1)
	}

	if err := w.Update(); err != nil {
		fmt.Println("Error getting weather data (", err, ")")
		os.Exit(-1)
	}

	fmt.Println("Time: ", time.Unix(w.Current.Dt, 0))
	fmt.Println("Sun rise: ", time.Unix(w.Current.SunRise, 0))
	fmt.Println("Sun set: ", time.Unix(w.Current.SunSet, 0))

	fmt.Println("Temperature: ", w.Current.Temperature, "°C")
	fmt.Println("Humindity: ", w.Current.Humidity, "%")
	fmt.Println("Pressure: ", w.Current.Pressure, "mbar")
	fmt.Println("Rain: ", w.Current.Rain.OneHour, "mm")
	fmt.Println("Wind speed: ", w.Current.WindSpeed, "m/s")
	fmt.Println("Wind Degree: ", w.Current.WindDegree, "°")

	fmt.Println(w)
}
