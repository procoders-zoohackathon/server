package main

import (
	"errors"
)

func sendMessage(id int, message string) error {
	clientListMu.RLock()
	if id >= len(clientList) || id < -1 {
		return errors.New("invalid id")
	}
	so := clientList[id]
	clientListMu.RUnlock()

	so.EmitJSON(struct {
		Request string `json:"request"`
		Data    string `json:"data"`
	}{"message", message})
	return nil
}
