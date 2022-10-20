package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func makeURL(host, port, methodPath string) string {
	return "http://" + host + ":" + port + methodPath
}

func GetHttpResponse(host, port, methodPath string) ([]byte, error) {
	client := http.Client{Timeout: time.Duration(2) * time.Second}
	resp, err := client.Get(makeURL(host, port, methodPath))
	if err != nil {
		return nil, fmt.Errorf("httpclient - GetHttpResponse: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusInternalServerError {
			log.Println("ERROR:", host, ":", port, "/", methodPath, " - Response Status: ", resp.StatusCode)
		}
		return nil, fmt.Errorf("incorrect status code: %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("httpclient - GetHttpResponse: %w", err)
	}
	return body, nil
}
