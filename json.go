package main

import (
	"encoding/json"
	"time"
)

type JsonOutput struct {
	Time  time.Time `json:"dateUpdated"`
	Drips []Drip    `json:"drips"`
}

func marshallJson(d []Drip) ([]byte, error) {
	data := JsonOutput{
		Time:  time.Now(),
		Drips: d,
	}
	return json.Marshal(data)
}
