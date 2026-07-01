package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Location struct {
	lat, long float32
}

func validate_and_clean(coord string) (float32, error) {
	if coord == "" {
		return float32(0), nil
	}
	coord_clean, err := strconv.ParseFloat(coord, 7)
	return float32(coord_clean), err
}
func (loc Location) get_location() {
	var lat_in, long_in string
	fmt.Print("latitude (default is 0): ")
	fmt.Scanln(&lat_in)
	lat, err := validate_and_clean(lat_in)
	if err != nil {
		fmt.Println("invalid coordinates!")
		return
	}
	fmt.Println("selected latitude: ", lat)
	fmt.Print("longitude (default is 0): ")
	fmt.Scanln(&long_in)
	long, err := validate_and_clean(long_in)
	if err != nil {
		fmt.Println("invalid coordinates!")
		return
	}
	fmt.Println("selected latitude: ", long)
	loc.set_lat(lat)
	loc.set_long(long)
}
func (loc Location) set_lat(lat float32)   { loc.lat = lat }
func (loc Location) set_long(long float32) { loc.long = long }

func show_menu() {
	fmt.Println("--------------")
	fmt.Println("1. Input coordinates")
	fmt.Println("2. Show weather forecast")
	fmt.Println("3. Exit")
	fmt.Print("Option chosen: ")
}

type Forecast struct {
	Forecast_time    []string  `json:"time"`
	Temperature_2m   []float32 `json:"temperature_2m"`
	Precipitation    []float32 `json:"precipitation"`
	Surface_pressure []float32 `json:"surface_pressure"`
	Wind_speed_180m  []float32 `json:"wind_speed_180m"`
}

type Response struct {
	Hourly Forecast `json:"hourly"`
}

func (forecast Forecast) show_forecast() {
	fmt.Printf("Time\t\t\tTemperature (C)\tPrecipitation (mm^2)\tPressure (hPa)\tWind (km/h)\n")
	for i := range forecast.Forecast_time {
		fmt.Printf("%s\t%.1f\t\t%.f\t\t\t%.f\t\t%.f\n", forecast.Forecast_time[i], forecast.Temperature_2m[i], forecast.Precipitation[i], forecast.Surface_pressure[i], forecast.Wind_speed_180m[i])
	}
}

func get_weather(location Location) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation,surface_pressure,wind_speed_180m&forecast_days=1", location.lat, location.long))
	if err != nil {
		fmt.Println("couldnt get weaTHER DATA!")
		return
	}
	defer resp.Body.Close()
	var data = Response{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	data.Hourly.show_forecast()
}

func main() {
	fmt.Println("WEATHER FO TODAY")
	location := Location{}
	nav := 0
	for ok := true; ok; ok = (nav != 3) {
		show_menu()
		fmt.Scanln(&nav)
		switch nav {
		case 1:
			location.get_location()
			break
		case 2:
			get_weather(location)
			break
		case 3:
			return
		default:
			break
		}
	}
}
