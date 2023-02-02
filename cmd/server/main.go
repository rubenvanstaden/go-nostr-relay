package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/rubenvanstaden/go-nostr-relay/core"
	"github.com/rubenvanstaden/go-nostr-relay/inmem"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type socket struct {
	hub             *Hub
	eventRepository core.EventRepository
}

// Handle websocket requests from the peer.
func (s *socket) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Upgrade the http protocol to a websocket.
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new client abstraction for every connection to relay.
	client := newClient(s.hub, c, s.eventRepository)

	// Register and connect newly created client to hub.
	client.hub.register <- client

	// Instantiate a write and read goroutines for every client connection established.
	go client.write()
	go client.read()
}

func main() {

	flag.Parse()
	log.SetFlags(0)

	hub := newHub()
	go hub.run()

	repository := inmem.New()

	h := &socket{
		hub:             hub,
		eventRepository: repository,
	}

	http.Handle("/", h)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
