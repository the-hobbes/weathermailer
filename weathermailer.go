package main

/*
	./weathermailer \
		-appid $APPID \
		-password $PASSWORD \
		-destinations $DESTINATIONS
*/

import (
	"flag"
)

func GetFlags() (ApiInfo, ConnectionInfo, bool) {
	// api request flags
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

	// connection info flags
	sender := flag.String(
		"sender",
		"phelan.vendeville@gmail.com",
		"An email address representing the source of the mail")
	port := flag.String(
		"port", "465", "The port to use for the SMTP connection. Defaults to 465.")
	host := flag.String(
		"host", "smtp.gmail.com", "The sending SMTP server. Defaults to gmail.")
	password := flag.String(
		"password", "", "The password associated with the sender.")
	var destinationList DestinationAddresses
	flag.Var(
		&destinationList,
		"destinations",
		"A comma separated list of email addresses to send to.")

	// generator flags
	generate := flag.Bool(
		"generate_proto",
		false,
		"Whether or not to generate the proto and proto library.")

	flag.Parse()
	apiInfo := ApiInfo{}
	connInfo := ConnectionInfo{}

	if *generate == true {
		return apiInfo, connInfo, *generate
	}

	// set api info struct
	apiInfo.city = *city
	apiInfo.countryCode = *countryCode
	apiInfo.apiId = *apiId
	apiInfo.units = *units
	apiInfo.lines = *lines

	// set connection info struct
	connInfo.sender = *sender
	connInfo.port = *port
	connInfo.host = *host
	connInfo.password = *password
	connInfo.destinations = destinationList

	return apiInfo, connInfo, false
}

func main() {
	apiInfo, connInfo, generate := GetFlags()
	if generate {
		DoGenerateProto()
	} else {
		subject, body := DoForecast(&apiInfo)
		DoMail(&connInfo, subject, body)
	}
	// TODO: instead of commandline flags, add a config script
}
