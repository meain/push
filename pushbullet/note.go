package pushbullet

import (
	"bytes"
	"encoding/json"
)

type pushBody struct {
	Title string `json:"string"`
	Body  string `json:"body"`
	Type  string `json:"type"`
	Iden  string `json:"device_iden,omitempty"`
}

// PushNote pushes a note to pushbullet
func PushNote(title string, body string, token string, device string) error {
	// curl --header 'Access-Token: <your_access_token_here>' \
	//      --header 'Content-Type: application/json' \
	//      --data-binary '{"body":"Space Elevator, Mars Hyperloop, Space Model S (Model Space?)","title":"Space Travel Ideas","type":"note"}' \
	//      --request POST \
	//      https://api.pushbullet.com/v2/pushes

	userDevice, err := GetDefaultDevice(device, token)
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
