package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/rubenvanstaden/go-nostr-relay/core"
	"github.com/rubenvanstaden/go-nostr-relay/inmem"
	"github.com/rubenvanstaden/go-nostr-relay/relay"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type socket struct {
	relay core.RelayService
}

func (s *socket) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {

		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var tmp []json.RawMessage

		err = json.Unmarshal(message, &tmp)
		if err != nil {
			panic(err)
		}

		msgType, err := core.MessageFromBytes(tmp[0])

		switch msgType {
		case core.MessageEvent:

			log.Printf("Event: %s", msgType)

			var event core.Event
			err = json.Unmarshal(tmp[1], &event)
			if err != nil {
				panic(err)
			}

			s.relay.Publish(&event)

		default:
			panic(fmt.Errorf("unknown message type"))
		}

		log.Printf("recv: %s", message)

		err = c.WriteMessage(mt, []byte("dergigi"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {

	flag.Parse()
	log.SetFlags(0)

	repository := inmem.New()
	service := relay.New(repository)

	h := &socket{
		relay: service,
	}

	http.Handle("/", h)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
