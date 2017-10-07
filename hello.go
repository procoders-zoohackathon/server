package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/TF2Stadium/wsevent"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(_ *http.Request) bool { return true }}
	server   = wsevent.NewServer(JSONCodec{}, func(_ *wsevent.Client, _ struct{}) interface{} {
		return errors.New("no such request")
	})
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/index.html")
	})

	http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/index.js")
	})

	http.HandleFunc("/connect", SocketHandler)
	go http.ListenAndServe("0:8080", nil)
	server.OnDisconnect = func(string, *jwt.Token) {
		log.Println("client disconnected")
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		params := strings.Split(text, ",")
		if len(params) != 2 {
			log.Print("invalid message format")
			continue
		}
		id, err := strconv.Atoi(params[0])
		if err != nil {
			log.Print(err)
			continue
		}
		if err := sendMessage(id, params[1]); err != nil {
			log.Print(err)
		}
	}
}
