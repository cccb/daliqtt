package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Global state
var lights []int = []int{23, 42, 0, 254}

func parseCommand(query string) (int, int, error) {
	cmd, err := url.QueryUnescape(query)
	if err != nil {
		log.Println("An error occured:", err)
		return 0, 0, err
	}

	tokens := strings.Split(cmd, " ")
	if len(tokens) != 4 {
		return 0, 0, fmt.Errorf("Invalid command")
	}

	id, err := strconv.ParseInt(tokens[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	value, err := strconv.ParseInt(tokens[3], 10, 64)
	if err != nil {
		return int(id), 0, err
	}

	return int(id), int(value), nil
}

func handleRetrieve(req *http.Request) []byte {
	res := ""
	for _, v := range lights {
		res += fmt.Sprintf("%d ", v)
	}

	res += "\r\n"

	return []byte(strings.TrimSpace(res))
}

func handleUpdate(req *http.Request) []byte {
	id, value, err := parseCommand(req.URL.RawQuery)
	if err != nil {
		return []byte(err.Error())
	}

	if id >= len(lights) {
		return []byte("ERROR")
	}

	log.Println("Updating", id, "set to", value)
	lights[id] = value

	return []byte("")
}

func main() {
	fmt.Println("Starting licht.cgi stub server")

	_ = lights

	http.HandleFunc(
		"/cgi-bin/licht.cgi",
		func(res http.ResponseWriter,
			req *http.Request) {
			if req.Method == "POST" {
				res.Write(handleUpdate(req))
			} else {
				res.Write(handleRetrieve(req))
			}
		})

	log.Println("Listening to: localhost:2299")
	http.ListenAndServe("localhost:2299", nil)
}
