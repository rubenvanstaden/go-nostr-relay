package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/rubenvanstaden/go-nostr-relay/core"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func socket(w http.ResponseWriter, r *http.Request) {

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

        var event core.Event

        switch msgType {
        case core.MessageEvent:
            log.Printf("Event: %s", msgType)
            err = json.Unmarshal(tmp[1], &event)
            if err != nil {
                panic(err)
            }
        default:
            panic(fmt.Errorf("unknown message type"))
        }

        log.Println(string(tmp[0]))
        log.Println(event)

		log.Printf("recv: %s", message)

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {

	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/nostr", socket)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
