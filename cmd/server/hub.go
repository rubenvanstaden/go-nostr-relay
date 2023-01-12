package main

import "github.com/rubenvanstaden/go-nostr-relay/core"

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {

	// Registered clients.
	clients map[*Client]struct{}

	// Inbound events from the clients that should be broadcasted to subscribed clients.
	broadcast chan *core.Event

	// Connected clients. Register requests from the clients.
	register chan *Client

	// Disconnect clients. Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *core.Event),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]struct{}),
	}
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
		case event := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- event:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
