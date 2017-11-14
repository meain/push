package pushbullet

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
)

// Device contains information about one user device
type Device struct {
	Iden  string `json:"iden"`
	Nick  string `json:"nickname"`
	Token string `json:"push_token"`
}

// GetDevices provides a method to retrieve device information from pushbullet
func GetDevices(token string) []Device {
	// curl --header 'Access-Token: <your_access_token_here>' \
	//      https://api.pushbullet.com/v2/devices

	pbURL := "https://api.pushbullet.com"
	headers := []header{
		{"Access-Token", token},
	}

	res, err := makeRequest(pbURL+"/v2/devices", "GET", headers, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var userDevices deviceResponse
	json.NewDecoder(res.Body).Decode(&userDevices)
	return userDevices.Devices
}

// GetDefaultDevice return data about the device user would like to use
func GetDefaultDevice(iden string, token string) (Device, error) {
	var userDevice Device
	devices := GetDevices(token)
	for _, v := range devices {
		if strings.Compare(v.Iden, iden) == 0 {
			userDevice = v
		}
	}
	if (Device{}) == userDevice {
		return Device{}, errors.New("User device not found")
	}
	return userDevice, nil
}
