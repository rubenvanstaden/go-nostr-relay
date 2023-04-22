package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

// Alice pushes events to the relay that will be consumed by Bob at every ticker.

var   addr = flag.String("addr", "localhost:8080", "http service address")

var e1 = `{"id":"1","kind":"1","content":"hello world 1"}`
var e2 = `{"id":"2","kind":"1","content":"hello world 2"}`
var e3 = `{"id":"3","kind":"1","content":"hello world 3"}`

func main() {

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

:q
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for _, event := range []string{e1, e2, e3} {
		body := fmt.Sprintf(`["EVENT", %s]`, event)
		var data = []byte(body)
		err := c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
