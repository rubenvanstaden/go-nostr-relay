package main

import (
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
	send chan *core.Event

	// A set of subscribed filters to be applied before adding events to the outbound channel.
	subscriptions map[core.SubId]*core.Filter

	// The client is responsible to adding events to the repository
	events core.EventRepository
}

func newClient(hub *Hub, conn *websocket.Conn, repository core.EventRepository) *Client {
	return &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan *core.Event, 100),
		subscriptions: make(map[core.SubId]*core.Filter),
		events:        repository,
	}
}

func (s *Client) read() {

	// Unregister and close connection if reading a message fails.
	defer func() {
		s.hub.unregister <- s
		s.conn.Close()
	}()

	for {

		msgType, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		log.Printf("read msg type: %d", msgType)

		msg := core.DecodeClientMessage(message)

		switch msg.Type {
		case core.MessageEvent:
			if msg.Event.Kind == 0 {
				s.setMetadata(msg.Event)
			} else if msg.Event.Kind == 1 {
				s.textNote(msg.Event)
			} else if msg.Event.Kind == 2 {
				s.recommendServer(msg.Event)
			}
		case core.MessageRequest:
			fmt.Println("not implemented")
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

			// TODO
			// Get subscription that allows event to be broadcasted to client.
			// Send a relay notice if errors occur.

			e := core.RelayEvent{
				SubId: core.SubId("123"),
				Event: event,
			}

			err := s.conn.WriteMessage(websocket.TextMessage, e.Encode())
			if err != nil {
				log.Println("ERROR write:", err)
				return
			}
		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := s.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}
	}
}

func (s *Client) setMetadata(event *core.Event) error {
	op := "client.SetMetadata"
	log.Fatalf("%s: not implemented", op)
	return nil
}

func (s *Client) textNote(event *core.Event) error {

	// Store event to be pulled and filtered by future subscribers.
	s.events.Store(event)

	// After the event is stored broadcast t to all registered clients.
	s.hub.broadcast <- event

	return nil
}

func (s *Client) recommendServer(event *core.Event) error {
	op := "client.RecommendServer"
	log.Fatalf("%s: not implemented", op)
	return nil
}
