package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const IP_API = `https://httpbin.org/ip`

type IP struct {
	Origin *string `json:"origin"`
}

func ipGet() (*string, error) {
	resp, err := http.Get(IP_API)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ip := IP{}
	if err := json.Unmarshal(body, &ip); err != nil {
		return nil, err
	}
	return ip.Origin, nil
}
