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

// func (s *socket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//
// 	c, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Print("upgrade:", err)
// 		return
// 	}
// 	defer c.Close()
//
// 	for {
//
// 		mt, message, err := c.ReadMessage()
// 		if err != nil {
// 			log.Println("read:", err)
// 			break
// 		}
//
// 		var tmp []json.RawMessage
//
// 		err = json.Unmarshal(message, &tmp)
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		msgType, err := core.MessageFromBytes(tmp[0])
//
// 		switch msgType {
// 		case core.MessageEvent:
//
// 			log.Printf("Event: %s", msgType)
//
// 			var event core.Event
// 			err = json.Unmarshal(tmp[1], &event)
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			s.relay.Publish(&event)
//
// 		default:
// 			panic(fmt.Errorf("unknown message type"))
// 		}
//
// 		log.Printf("recv: %s", message)
//
// 		err = c.WriteMessage(mt, []byte("dergigi"))
// 		if err != nil {
// 			log.Println("write:", err)
// 			break
// 		}
// 	}
// }

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
