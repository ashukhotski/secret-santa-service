// util.go
package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendSlackMessage(url string, resType string, msg string) error {
	jsonData := new(bytes.Buffer)
	json.NewEncoder(jsonData).Encode(&SlackMessage{resType, msg})
	req, err := http.NewRequest("POST", url, jsonData)
	//var jsonData = []byte(`{"response_type": "` + resType + `", "text": "` + msg + `"}`)
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}

func String(s string) *string {
	return &s
}
