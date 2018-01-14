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
  description  string
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

func CreateFolksySaying(p *ParsedApiResponse) string {
  weather := p.Weather[0].Main  // TODO: get the most frequent category
  // loop through each of the weather messages in the proto
  // if the type matches the weather,
  // then grab a random saying from the repeated sayings field.
  // if no match is found, return
  // "Pretty weird, 'cause we don't have a folksy saying for that particular
  // kinda weather!"
  return saying
}

func CreateMessage(
  a *apiInfo, p *ParsedApiResponse, avg, saying string) string {
	// create the body and subject of the email that will be sent
  short_desc := p.Weather[0].Main
  long_desc := p.Weather[0].Description
  city := a.city
  subject := fmt.Sprintf("Today's weather is: %s", short_desc)
  body := fmt.Sprintf("Today in %s, the average temperature will be %d. " +
                      "Expect %d.\nIn other words, it'll be... %s",
                      city, avg, long_desc, saying)

  return body
}

func main() {
	apiInfo := SetForecastFlags()
	response := MakeOpenWeatheRequest(&apiInfo)
	parsed := ParseOpenWeatherResponse(response)
	forcastedAverage := ComputeForecastedAverage(&parsed)
  saying := CreateFolksySaying(&parsed)
	message := CreateMessage(&apiInfo, &parsed, forcastedAverage, saying)

	// TODO: make the call to mail...
}
