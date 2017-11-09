package transfersh

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

// TransferFile uploads file to transfer.sh and returns the link of file
func TransferFile(filePath string) string {
	filePath = expandFilePath(filePath)
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
