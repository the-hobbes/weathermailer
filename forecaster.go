package main

import (
  "fmt"
  "flag"
  "io/ioutil"
  "log"
  "net/http"
  // "strings"
)

type ApiInfo struct {
  city 				string
  countryCode string
  apiId  			string
  units     	string
  lines				string
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

func MakeRequest(a *ApiInfo) string {
	url := fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/forecast?q=%s,%s&APPID=%s&units=%s&cnt=%s", 
		a.city, 
		a.countryCode, 
		a.apiId, 
		a.units, 
		a.lines)
	
	log.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	bodyString := string(body)
	return bodyString
	// TODO: Next up, marshall the response into JSON.
}

func main() {
	apiInfo := SetForecastFlags()
	response := MakeRequest(&apiInfo)
	fmt.Println(response)
}