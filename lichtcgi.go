package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func (self *LichtCgi) FetchLights() ([]Light, error) {
	res, err := http.Get(self.Url + "/cgi-bin/licht.cgi")
	if err != nil {
		return []Light{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Light{}, err
	}

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
