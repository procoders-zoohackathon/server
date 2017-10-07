package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/TF2Stadium/wsevent"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/procoders-zoohackathon/server/reader"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(_ *http.Request) bool { return true }}
	server   = wsevent.NewServer(JSONCodec{}, func(_ *wsevent.Client, _ struct{}) interface{} {
		return errors.New("no such request")
	})
)

var (
	addr = flag.String("addr", ":8080", "server address")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/index.html")
	})

	http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/index.js")
	})

	http.HandleFunc("/connect", SocketHandler)
	log.Printf("hosting on %s", *addr)

	go func() {
		if err := http.ListenAndServe(*addr, nil); err != nil {
			log.Fatal(err)
		}
	}()
	server.OnDisconnect = func(string, *jwt.Token) {
		log.Println("client disconnected")
	}

	for {
		rdr := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := rdr.ReadString('\n')
		values, err := csv.NewReader(bytes.NewReader([]byte(text))).Read()
		id, err := strconv.Atoi(values[0])
		if err != nil {
			log.Print(err)
			continue
		}

		alert, err := reader.NewAlert(values[1:])
		if err != nil {
			log.Print(err)
			continue
		}
		if err := sendMessage(id, *alert); err != nil {
			log.Print(err)
		}
	}
}
