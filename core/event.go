package core

// https://github.com/nostr-protocol/nips/blob/master/01.md

import (
	"encoding/json"
	"strings"
)

type MessageType uint64

const (
	MessageUnknown MessageType = iota + 1
	MessageEvent
	MessageRequest
	MessageClose
)

func (s MessageType) String() string {
	switch s {
	case MessageEvent:
		return "EVENT"
	case MessageRequest:
		return "REQ"
	case MessageClose:
		return "CLOSE"
	default:
		return "UNKOWN"
	}
}

func DecodeMessageType(data []byte) MessageType {

	s := strings.Trim(string(data), "\"")

	switch s {
	case "EVENT":
		return MessageEvent
	case "REQ":
		return MessageRequest
	case "CLOSE":
		return MessageClose
	default:
		return MessageUnknown
	}
}

type Tag []string

type EventId string

type Event struct {
	Id      EventId `json:"id"`
	Kind    uint8   `json:"kind,string"`
	Content string  `json:"content"`
	// Pubkey    string  `json:"pubkey"`
	// CreatedAt string  `json:"created_at"`
	// Tags      []Tag   `json:"tags"`
	// Sig       string  `json:"sig"`
}

func (s Event) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

type SubId string

func (s SubId) String() string {
	return string(s)
}
