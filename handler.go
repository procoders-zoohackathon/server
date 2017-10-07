package main

import (
	"github.com/TF2Stadium/wsevent"
	"log"
	"net/http"
	"sync"
)

var (
	clientList   []*wsevent.Client
	clientListMu = new(sync.RWMutex)
)

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	so, err := server.NewClient(upgrader, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clientListMu.Lock()
	log.Printf("new client %d created", len(clientList))
	clientList = append(clientList, so)
	clientListMu.Unlock()
}
