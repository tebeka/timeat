// package timeat finds current time at address using Google Maps API

package timeat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"strings"
	"time"
)

const (
	apiBase = "https://maps.googleapis.com/maps/api"
	geoURL  = apiBase + "/geocode/json"
	tzURL   = apiBase + "/timezone/json"

	// Version is the package version
	Version = "2.1.0"
)

var apiKey = ""

type apiLoc struct {
	Address  string `json:"formatted_address"`
	Geomerty struct {
		Location struct {
			Lat float32 `json:"lat"`
			Lng float32 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
}

type apiGEOReply struct {
	Status    string   `json:"status"`
	Locations []apiLoc `json:"results"`
}

type apiTZReply struct {
	DST    float32 `json:"dstOffset"`
	Offset float32 `json:"rawOffset"`
	Status string  `json:"status"`
}

// TimeInfo holds local time at an address
type TimeInfo struct {
	Address string    // Address
	Time    time.Time // Local time
}

func findAPIKey() string {
	key := os.Getenv("TIMEIT_API_KEY")
	if key != "" {
		return key
	}

	user, err := user.Current()
	if err != nil {
		return ""
	}

	keyFile := path.Join(user.HomeDir, ".config", "timeat", "api-key.txt")
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

// SetAPIKey sets the API key
// If key is "" will try to read from TIMEIT_API_KEY environment variable or
// from ~/.config/timeat/api-key.txt
func SetAPIKey(key string) {
	if key == "" {
		key = findAPIKey()
	}

	if key != "" {
		apiKey = key
	}
}

// apiCall calls Google Geo API and populates reply
func apiCall(url string, vals url.Values, reply interface{}) error {
	if apiKey != "" && vals.Get("key") == "" {
		vals.Add("key", apiKey)
	}
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
		return fmt.Errorf("can't decode reply - %s", err)
	}

	return nil
}

// TimeAt return list of location matching address and the local time in each
// of them
func TimeAt(address string) ([]TimeInfo, error) {

	var times []TimeInfo

	vals := url.Values{}
	vals.Add("address", address)
	grep := apiGEOReply{}
	if err := apiCall(geoURL, vals, &grep); err != nil {
		return nil, fmt.Errorf("can't get geo for %s - %s", address, err)
	}

	if grep.Status != "OK" {
		return nil, fmt.Errorf("bad geo status - %s", grep.Status)
	}

	now, err := NTPTime()
	if err != nil {
		now = time.Now().UTC()
	}

	for _, loc := range grep.Locations {
		vals = url.Values{}
		vals.Add("timestamp", fmt.Sprintf("%d", now.Unix()))
		geo := loc.Geomerty.Location
		vals.Add("location", fmt.Sprintf("%f,%f", geo.Lat, geo.Lng))
		tz := apiTZReply{}
		if err := apiCall(tzURL, vals, &tz); err != nil {
			return nil, fmt.Errorf("can't get geo for %s - %s", address, err)
		}
		if tz.Status != "OK" {
			return nil, fmt.Errorf("bad tz status - %s", tz.Status)
		}
		dst, offset := time.Duration(tz.DST), time.Duration(tz.Offset)
		local := now.Add((dst + offset) * time.Second)

		times = append(times, TimeInfo{loc.Address, local})
	}

	return times, nil
}
