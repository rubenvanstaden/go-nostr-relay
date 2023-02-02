package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/rubenvanstaden/go-nostr-relay/core"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {

	// Required to put messages from websocket into hub.
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound event messages.
	send chan *core.RelayEvent

	// A set of subscribed filters to be applied before adding events to the outbound channel.
	subscriptions map[core.SubId]*core.Filter

	// The client is responsible to adding events to the repository
	events core.EventRepository
}

func newClient(hub *Hub, conn *websocket.Conn, repository core.EventRepository) *Client {
	return &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan *core.RelayEvent, 100),
		subscriptions: make(map[core.SubId]*core.Filter),
		events:        repository,
	}
}

func (s *Client) validSubId(id core.SubId) error {

	if subId, found := s.subscriptions[id]; found {
		return core.Errorf(core.ErrorConflict, "Subscription ID already stored: %s", subId)
	}

	return nil
}

func (s *Client) subscribed(e *core.Event) *core.RelayEvent {

	for id, filter := range s.subscriptions {
		for _, eid := range filter.Ids {
			if eid == e.Id {
				return &core.RelayEvent{
					SubId: id,
					Event: e,
				}
			}
		}
	}

	return nil
}

func (s *Client) read() {

	// Unregister and close connection if reading a message fails.
	defer func() {
		s.hub.unregister <- s
		s.conn.Close()
	}()

	for {

		// If wsType=1 then its a text message.
		// TODO: Implement support with close message type.
		wsType, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		log.Printf("read msg type: %d", wsType)

		//msg := core.DecodeClientMessage(message)

		var tmp []json.RawMessage

		err = json.Unmarshal(message, &tmp)
		if err != nil {
			log.Fatalln("unable to unmarshal client msg")
		}

		// Set message type from first array item.
		msgType := core.DecodeMessageType(tmp[0])

		switch msgType {
		case core.MessageEvent:

			var event core.Event
			err = json.Unmarshal(tmp[1], &event)
			if err != nil {
				panic(err)
			}

			if event.Kind == 0 {
				s.setMetadata(&event)
			} else if event.Kind == 1 {
				s.textNote(&event)
			} else if event.Kind == 2 {
				s.recommendServer(&event)
			}
		case core.MessageRequest:

			var subId core.SubId
			err = json.Unmarshal(tmp[1], &subId)
			if err != nil {
				panic(err)
			}

			err = s.validSubId(subId)
			if err != nil {
				panic(err)
			}

			var filters []*core.Filter
			err = json.Unmarshal(tmp[2], &filters)
			if err != nil {
				panic(err)
			}

            err := s.request(subId, filters)
			if err != nil {
				panic(err)
			}
		case core.MessageClose:
			fmt.Println("not implemented")
		default:
			panic(fmt.Errorf("unknown client message type"))
		}
	}

}

func (s *Client) write() {

	ticker := time.NewTicker(time.Second)
	defer func() {
		ticker.Stop()
		s.conn.Close()
	}()

	for {
		select {
		case event, ok := <-s.send:

			s.conn.SetWriteDeadline(time.Now().Add(writeWait))

			// The hub closed the channel.
			if !ok {
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := s.conn.WriteMessage(websocket.TextMessage, event.Encode())
			if err != nil {
				log.Println("ERROR write:", err)
				return
			}

			// TODO
			// Send a relay notice if errors occur.

		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := s.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}
	}
}

func (s *Client) request(subId core.SubId, filters []*core.Filter) error {

    log.Println("requesting")

	for _, filter := range filters {

		for _, id := range filter.Ids {

            event, err := s.events.FindById(id)
            if err != nil {
                return err
            }

            e := &core.RelayEvent{
                SubId: subId,
                Event: event,
            }

            s.send <- e
        }
	}

	return nil
}

func (s *Client) setMetadata(event *core.Event) error {
	op := "client.SetMetadata"
	log.Fatalf("%s: not implemented", op)
	return nil
}

func (s *Client) textNote(event *core.Event) error {

	// Store event to be pulled and filtered by future subscribers.
	err := s.events.Store(event)
	if err != nil {
		panic(err)
	}

	// After the event is stored broadcast it to all registered clients.
	s.hub.Broadcast(event)

	return nil
}

func (s *Client) recommendServer(event *core.Event) error {
	op := "client.RecommendServer"
	log.Fatalf("%s: not implemented", op)
	return nil
}
