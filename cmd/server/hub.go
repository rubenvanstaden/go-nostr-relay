package main

import "github.com/rubenvanstaden/go-nostr-relay/core"

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {

	// Registered clients.
	clients map[*Client]struct{}

	// Inbound events from the clients that should be broadcasted to subscribed clients.
	events chan *core.Event

	// Connected clients. Register requests from the clients.
	register chan *Client

	// Disconnect clients. Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		events:     make(chan *core.Event),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]struct{}),
	}
}

func (s *Hub) Broadcast(event *core.Event) {
	s.events <- event
}

// Listens to registered and unregistered client and incoming messages.
// Incoming messages will be distributed to registered clients.
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case event := <-h.events:
			for client := range h.clients {

				relayEvent := client.subscribed(event)

				if relayEvent != nil {
					select {
					case client.send <- relayEvent:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
