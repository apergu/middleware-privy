package helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func HttpRequest(method, url string, data []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)

	log.Println("resp : ", resp)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

func IFEmpty(input, defaultValue string) string {
	if input == "" {
		return defaultValue
	}
	return input
}
