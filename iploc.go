package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Location struct {
	Country   string `json:"country"`
	City      string `json:"city"`
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

func getIPLoc(ipAddress string) (*Location, error) {
	apiURL := fmt.Sprintf("https://sys.airtel.lv/ip2country/%s/?full=true", ipAddress)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, errors.New("can't talk with IP information server")
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) == "wrong_ip" {
		return nil, errors.New("ip address is malformed")
	}

	location := Location{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func logLocation(ipAddress string, location Location) {
	fmt.Printf("ip %s is located in %s (%s) at %s, %s \n", ipAddress, location.City, location.Country, location.Latitude, location.Longitude)
}

func main() {
	binName := os.Args[0]
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage : %s <ip address>\n", binName)
		return
	}

	ipAddress := os.Args[1]
	location, err := getIPLoc(ipAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s : %s\n", binName, err)
		os.Exit(1)
	}

	logLocation(ipAddress, *location)
}
