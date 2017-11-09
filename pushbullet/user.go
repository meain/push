package pushbullet

import (
	"encoding/json"
	"log"
)

// User has basic user info about the pushbullet user
type User struct {
	Iden            string `json:"iden"`
	Email           string `json:"email"`
	EmailNormalized string `json:"email_normalized"`
	Name            string `json:"name"`
	ImageURL        string `json:"image_url"`
}

// GetUser retrieves user data from pushbullet
func GetUser(token string) User {
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

	var user User
	json.NewDecoder(res.Body).Decode(&user)
	return user
}
