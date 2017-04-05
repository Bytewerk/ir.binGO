package main

import (
	"net/http"
	"io/ioutil"
)

func SendCommand(path, value string) (string, error) {
	resp, err := http.Get(path + "?" + value)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bdy := string(body[:])
	return bdy, err
}