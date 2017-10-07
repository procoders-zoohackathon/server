package main

import (
	"encoding/json"
)

type JSONCodec struct{}

func (JSONCodec) ReadName(data []byte) string {
	var body struct {
		Request string
	}
	json.Unmarshal(data, &body)
	return body.Request
}

func (JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (JSONCodec) Error(err error) interface{} {
	return struct {
		Message string `json:"message"`
	}{err.Error()}
}
