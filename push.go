package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type userInfo struct {
	Iden            string `json:"iden"`
	Email           string `json:"email"`
	EmailNormalized string `json:"email_normalized"`
	Name            string `json:"name"`
	ImageURL        string `json:"image_url"`
}

type device struct {
	Iden  string `json:"iden"`
	Nick  string `json:"nickname"`
	Token string `json:"push_token"`
}

type deviceResponse struct {
	Devices []device `json:"devices"`
}

type pushBody struct {
	Title string `json:"string"`
	Body  string `json:"body"`
	Type  string `json:"type"`
	Iden  string `json:"device_iden,omitempty"`
}

type header struct {
	name  string
	value string
}

type confData struct {
	Token  string `yaml:"token"`
	Device string `yaml:"device"`
}

func expandFilePath(filePath string) string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	return strings.Replace(filePath, "~", homeDir, 1)
}

func (c *confData) getConf(file string) *confData {
	file = expandFilePath("~/.push.yaml")
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func makeRequest(url string, reqType string, headers []header, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(reqType, url, body)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range headers {
		req.Header.Set(v.name, v.value)
	}

	res, err := client.Do(req)
	return res, err
}

func getUser(token string) userInfo {
	// curl --header 'Access-Token: <your_access_token_here>' \
	//      https://api.pushbullet.com/v2/users/me

	pbURL := "https://api.pushbullet.com"
	headers := []header{
		{"Access-Token", token},
	}

	res, err := makeRequest(pbURL+"/v2/users/me", "GET", headers, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var userData userInfo
	json.NewDecoder(res.Body).Decode(&userData)
	return userData
}

func getDevices(token string) []device {
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

func pushNote(title string, body string, token string, device string) error {
	// curl --header 'Access-Token: <your_access_token_here>' \
	//      --header 'Content-Type: application/json' \
	//      --data-binary '{"body":"Space Elevator, Mars Hyperloop, Space Model S (Model Space?)","title":"Space Travel Ideas","type":"note"}' \
	//      --request POST \
	//      https://api.pushbullet.com/v2/pushes

	userDevice, err := getUserDevice(device, token)
	if err != nil {
		return err
	}

	pbURL := "https://api.pushbullet.com"
	headers := []header{
		{"Access-Token", token},
		{"Content-Type", "application/json"},
	}
	content := pushBody{
		title,
		body,
		"note",
		userDevice.Iden,
	}

	reqContent, err := json.Marshal(content)
	res, err := makeRequest(pbURL+"/v2/pushes", "POST", headers, bytes.NewBuffer(reqContent))
	defer res.Body.Close()
	return err
}

func getUserDevice(deviceName string, token string) (device, error) {
	var userDevice device
	devices := getDevices(token)
	for _, v := range devices {
		if strings.Compare(v.Nick, deviceName) == 0 {
			userDevice = v
		}
	}
	if (device{}) == userDevice {
		return device{}, errors.New("User device not found")
	}
	return userDevice, nil
}

// Uploads file to transfer.sh and returns the link of file
func transferFile(filePath string) string {
	filePath = expandFilePath(expandFilePath(expandFilePath(filePath)))
	fileName := filepath.Base(filePath)
	transferURL := "https://transfer.sh/" + fileName

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	request, err := http.NewRequest("PUT", transferURL, bytes.NewReader(fileData))
	response, err := client.Do(request)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(responseBody)
}

func transferPushFile(token string, device string, fileName string) error {
	s := transferFile(fileName)
	fmt.Println(s)

	err := pushNote("Message", s, token, device)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var conf confData
	conf.getConf(expandFilePath("~/.push.yaml"))

	if len(os.Args) != 3 {
		fmt.Println(`Usage: push [command] ...
Available commands:
	note - sens a note to your pushbullet device
		 - eg: 'push note "Get Schwifty"'
	file - uploads a file to transfer.sh and sends link to pushbullet device
		 - eg 'push file "Assignment.docx"'`)
		return
	}

	switch os.Args[1] {
	case "note":
		fmt.Println("Pushing note: ", os.Args[2])
		err := pushNote("Message", os.Args[2], conf.Token, conf.Device)
		if err != nil {
			log.Fatal(err)
		}
	case "file":
		fmt.Println("Pushing file: ", os.Args[2])
		err := transferPushFile(conf.Token, conf.Device, os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println(`Usage: push [command] ...
Available commands:
	note - sens a note to your pushbullet device
		 - eg: 'push note "Get Schwifty"'
	file - uploads a file to transfer.sh and sends link to pushbullet device
		 - eg 'push file "Assignment.docx"'`)
	}
}
