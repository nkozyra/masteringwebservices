package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

const HEADER = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"

type Weather struct {
	Main weatherItem
}

type weatherItem struct {
	Temp float64
}

type Geo struct {
	Latitude  float64
	Longitude float64
}

type GeoIP struct {
	IP        string  "json:ip"
	Latitude  float64 "json:latitude"
	Longitude float64 "json:longitude"
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Soap    SOAPBody
}

type SOAPBody struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Response SOAPPayload
}

type SOAPPayload struct {
	Temperature float64
}

func convertTemperature(temp float64, destination string) float64 {

	fahrenheit := (temp * (9 / 5)) - 459.67
	celsius := (fahrenheit - 32) * (5 / 9)

	if destination == "f" {

		return float64(fahrenheit)

	} else if destination == "c" {

		return celsius

	} else {
		return temp
	}

}

func getGeo(ip string) Geo {
	fmt.Println(ip)
	req := "http://freegeoip.net/json/65.32.98.14"

	res, _ := http.Get(req)
	jsonRes, _ := ioutil.ReadAll(res.Body)

	userIP := GeoIP{}
	err := json.Unmarshal(jsonRes, &userIP)
	if err != nil {
		fmt.Println(err)
	}

	g := Geo{Latitude: userIP.Latitude, Longitude: userIP.Longitude}

	return g
}

func getTemp(geoData Geo) float64 {

	req := "http://api.openweathermap.org/data/2.5/weather?lat=" + strconv.FormatFloat(geoData.Latitude, 'f', 3, 32) + "&lon=" + strconv.FormatFloat(geoData.Longitude, 'f', 3, 32)
	fmt.Println(req)
	res, _ := http.Get(req)
	jsonRes, _ := ioutil.ReadAll(res.Body)
	tmp := Weather{}
	err := json.Unmarshal(jsonRes, &tmp)
	if err != nil {
		fmt.Println(err)
	}
	return tmp.Main.Temp
}

func weatherRouter(w http.ResponseWriter, r *http.Request) {
	hostPort := (r.RemoteAddr)
	host, port, _ := net.SplitHostPort(hostPort)

	LatLon := getGeo(host)
	fmt.Println(LatLon)
	fmt.Println("Port", port)

	temperature := getTemp(LatLon)
	fmt.Println(temperature)

	soapy := SOAPEnvelope{}
	soapy.Soap = SOAPBody{}
	soapy.Soap.Response = SOAPPayload{Temperature: temperature}

	output, _ := xml.MarshalIndent(soapy, " ", "    ")
	fmt.Fprintln(w, HEADER)
	fmt.Fprintln(w, string(output))
}

func main() {

	fmt.Println("Starting API server")
	http.HandleFunc("/weather", weatherRouter)
	http.ListenAndServe(":8080", nil)

}
