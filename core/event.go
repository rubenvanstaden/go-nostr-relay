package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

// https://github.com/nostr-protocol/nips/blob/master/01.md

type MessageType uint8

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
	Kind    uint8   `json:"kind"`
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

type Filter struct {
}

type RelayNotice struct {
	Message string `json:"message"`
}

func (s RelayNotice) Encode() []byte {

	array := []string{"NOTICE", s.Message}

	bytes, err := json.Marshal(array)
	if err != nil {
		panic(err)
	}

	return bytes
}

type RelayEvent struct {
	SubId SubId  `json:"sub_id"`
	Event *Event `json:"event"`
}

func (s RelayEvent) Encode() []byte {

	array := []string{"EVENT", s.SubId.String(), s.Event.String()}

	bytes, err := json.Marshal(array)
	if err != nil {
		panic(err)
	}

	return bytes
}

type ClientMessage struct {
	Type   MessageType `json:"message_type"`
	SubId  SubId       `json:"sub_id"`
	Event  *Event      `json:"event"`
	Filter Filter      `json:"filter"`
}

func DecodeClientMessage(data []byte) *ClientMessage {

	var tmp []json.RawMessage

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		panic(err)
	}

	msg := &ClientMessage{}

	// Set message type from first array item.
	msg.Type = DecodeMessageType(tmp[0])

	// Set message data from second array item.
	switch msg.Type {
	case MessageEvent:
		err = json.Unmarshal(tmp[1], msg.Event)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("unknown message type"))
	}

	return msg
}

func EncodeRelayMessage() {

}
