package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ParsedApiResponse struct {
	// http://openweathermap.org/forecast5#JSON
	City struct {
		Name string `json:"name"`
	}
	List []struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		} `json:"Weather"`
	} `json:"list"`
}

type MessageBodyContents struct {
	city         string
	units        string
	averageTemp  string
	folksySaying string
	weather      string
}

type ApiInfo struct {
	city        string
	countryCode string
	apiId       string
	units       string
	lines       string
}

func SetForecastFlags() ApiInfo {
	// get city, ISO country code, APPID, units, and number of lines
	city := flag.String(
		"city",
		"vergennes",
		"The city for which to retrieve the forecast. Defaults to Vergennes VT.")
	countryCode := flag.String(
		"countrycode",
		"840",
		"The ISO country code of the city. Defaults to the United States.")
	apiId := flag.String(
		"appid",
		"",
		"The appid to use for the openweathermap API calls.")
	units := flag.String(
		"units",
		"imperial",
		"The temperature units to use. Defaults to imperial.")
	lines := flag.String(
		"lines",
		"8",
		"The number of lines to retrieve from the API. Defaults to 8.")
	flag.Parse()
	apiInfo := ApiInfo{}
	apiInfo.city = *city
	apiInfo.countryCode = *countryCode
	apiInfo.apiId = *apiId
	apiInfo.units = *units
	apiInfo.lines = *lines

	return apiInfo
}

func MakeOpenWeatheRequest(a *ApiInfo) []byte {
	url := fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/forecast?q=%s,%s&APPID=%s&units=%s&cnt=%s&mode=json",
		a.city,
		a.countryCode,
		a.apiId,
		a.units,
		a.lines)

	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	log.Println("openweathermap query completed successfully.")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	return body
}

func ParseOpenWeatherResponse(b []byte) ParsedApiResponse {
	// parse the response into a ParsedApiResponse
	var p ParsedApiResponse
	err := json.Unmarshal(b, &p)
	if err != nil {
		log.Panic(err)
	}
	log.Println("JSON successfully unmarshalled.")

	return p
}

func ComputeForecastedAverage(p *ParsedApiResponse) string {
	// calculate the average of a list of temperatures
	a := p.List
	sum := float64(0)
	for _, element := range a {
		sum += element.Main.Temp
	}
	avg := sum / float64(len(a))

	return strconv.FormatFloat(avg, 'f', -1, 32)
}

func CreateMessageBodyContents(
	a *apiInfo, temp string, saying, weather string) MessageBodyContents {
	msg := MessageBodyContents{}
	msg.city = a.city
	msg.units = a.units
	msg.averageTemp = temp
	msg.folksySaying = saying
	msg.weather = weather

	return msg
}

func SelectWeatherSaying(p *ParsedApiResponse) string {
	// logic to select the folksy saying based on weather
	// TODO
	return "S"
}

func CreateMessageBody(a *apiInfo, avg, saying string) string {
	// create the body of the email that will be sent
	// TODO, using CreateMessageBodyContents()
	return "S"
}

func main() {
	apiInfo := SetForecastFlags()
	response := MakeOpenWeatheRequest(&apiInfo)
	parsed := ParseOpenWeatherResponse(response)
	forcastedAverage := ComputeForecastedAverage(&parsed)
	saying := SelectWeatherSaying(&parsed)
	body := CreateMessageBody(&apiInfo, forcastedAverage, saying)

	// TODO: make the call to mail...
}
