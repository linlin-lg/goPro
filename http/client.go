package main

import (
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return b
}
