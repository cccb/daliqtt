package main

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

/*

Dali / licht.cgi interface:

Retrieve light values
GET http://dali/cgi-bin/licht.cgi

L1 L2 L3 L4

Ln = 0 .. 255

Set light value

POST http://dali/cgi-bin/licht.cgi?set N to V

N = 0 .. 4

V = 0 .. 255

*/

type LichtCgi struct {
	Url string
}

func NewLichtCgi(url string) *LichtCgi {
	cgi := &LichtCgi{
		Url: url,
	}

	return cgi
}

func (self *LichtCgi) FetchLights(retries int) ([]Light, error) {
	var (
		result []Light
		err    error
	)

	// Sometimes the server is acting strange.
	// Retry with some timeout until finally giving up
	for retry := 0; retry < retries; retry++ {
		result, err = self._fetchLights()
		if err == nil {
			break
		}

		log.Println("Retry after error while fetch state from server:", err)
		time.Sleep(1 * time.Second)
	}

	return result, err
}

func (self *LichtCgi) _fetchLights() ([]Light, error) {

	// As we seem to have issues with broken tcp connections,
	// let's create a fresh client.
	//
	// And create a new one for each request.
	// Even if the docs say you should reuse them.
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: tr}

	res, err := client.Get(self.Url + "/cgi-bin/licht.cgi")
	if err != nil {
		return []Light{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Light{}, err
	}
	log.Println("Received state:", string(body))

	lights := []Light{}
	values := strings.Split(strings.TrimSpace(string(body)), " ")
	for i, v := range values {
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return lights, err
		}
		lights = append(lights, Light{Id: i, Value: int(value)})
	}

	return lights, nil
}

func encodeCommand(id, value int) string {
	cmd := fmt.Sprintf("set %d to %d", id, value)

	return url.QueryEscape(cmd)
}

func (self *LichtCgi) Update(id, value int) error {
	if value > 255 {
		return fmt.Errorf("Value must be within 0 .. 255")
	}

	cmd := encodeCommand(id, value)
	_, err := http.PostForm(
		self.Url+"/cgi-bin/licht.cgi?"+cmd,
		url.Values{})

	return err
}
