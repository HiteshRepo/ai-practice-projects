package models

import "encoding/json"

type Weather struct {
	Temperature string `json:"temperature"`
	Unit        string `json:"unit"`
	Forecast    string `json:"forecast"`
}

func (w Weather) ToString() string {
	b, _ := json.Marshal(w)
	return string(b)
}

type Location struct {
	Address string `json:"address"`
}

func (l Location) ToString() string {
	return l.Address
}
