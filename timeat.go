/*
Print current time at address using Google Maps API

Example Usage

    $ ./timeat 'los angeles'
    Los Angeles, CA, USA: Fri May 23, 2014 22:33
    $ ./timeat paris
    Paris, France: Sat May 24, 2014 07:37
    Paris, TX, USA: Sat May 24, 2014 00:37
    Paris, TN 38242, USA: Sat May 24, 2014 00:37
    Paris, IL 61944, USA: Sat May 24, 2014 00:37
    Paris, KY 40361, USA: Sat May 24, 2014 01:37
    $
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	apiBase = "https://maps.googleapis.com/maps/api"
	geoUrl  = apiBase + "/geocode/json"
	tzUrl   = apiBase + "/timezone/json"

	Version = "0.1.1"
)

type Loc struct {
	Address  string `json:"formatted_address"`
	Geomerty struct {
		Location struct {
			Lat float32 `json:"lat"`
			Lng float32 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
}

type GEOReply struct {
	Status    string `json:"status"`
	Locations []Loc  `json:"results"`
}

type TZReply struct {
	DST    float32 `json:"dstOffset"`
	Offset float32 `json:"rawOffset"`
	Status string  `json:"status"`
}

// die prints error message and aborts the program
func die(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}

// apiCall calls Google Geo API and populates reply
func apiCall(url string, vals url.Values, reply interface{}) error {
	url = fmt.Sprintf("%s?%s", url, vals.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can't call %s - %s", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad return code - %s", resp.Status)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(reply); err != nil {
		return fmt.Errorf("error: can't decode reply - %s", err)
	}

	return nil
}

func main() {
	version := flag.Bool("version", false, "show version")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s ADDRESS\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		die("wrong number of arguments")
	}

	address := flag.Arg(0)

	vals := url.Values{}
	vals.Add("address", address)
	grep := GEOReply{}
	if err := apiCall(geoUrl, vals, &grep); err != nil {
		die("can't get geo for %s - %s", address, err)
	}

	if grep.Status != "OK" {
		die("error: bad status - %s", grep.Status)
	}

	if len(grep.Locations) == 0 {
		die("error: no locations found matching %s", address)
	}

	now := time.Now().UTC()
	for _, loc := range grep.Locations {
		vals = url.Values{}
		vals.Add("timestamp", fmt.Sprintf("%d", now.Unix()))
		geo := loc.Geomerty.Location
		vals.Add("location", fmt.Sprintf("%f,%f", geo.Lat, geo.Lng))
		tz := TZReply{}
		if err := apiCall(tzUrl, vals, &tz); err != nil {
			die("error: can't get geo for %s - %s", address, err)
		}
		if tz.Status != "OK" {
			die("error: bad status - %s", grep.Status)
		}
		dst, offset := time.Duration(tz.DST), time.Duration(tz.Offset)
		local := now.Add((dst + offset) * time.Second)
		fmt.Printf("%s: %s\n", loc.Address, local.Format("Mon Jan 2, 2006 15:04"))
	}
}
