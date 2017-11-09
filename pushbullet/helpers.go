package pushbullet

import (
	"io"
	"log"
	"net/http"
)

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
