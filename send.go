package main

import (
	"errors"
)

func sendMessage(id int, data interface{}) error {
	clientListMu.RLock()
	if id >= len(clientList) || id < -1 {
		return errors.New("invalid id")
	}
	so := clientList[id]
	clientListMu.RUnlock()

	so.EmitJSON(struct {
		Request string      `json:"request"`
		Data    interface{} `json:"data"`
	}{"message", data})
	return nil
}
