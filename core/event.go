package core

import (
	"fmt"
	"strings"
)

// https://github.com/nostr-protocol/nips/blob/master/01.md

// import "time"

// type Tag struct {
// 	Type string
// 	Key  []byte
// 	Url  string
// }

// type Event struct {
// 	Id      []byte
// 	Pubkey  []byte
// 	CreatAt time.Duration
// 	Kind    uint8
// 	Tags    []Tag
// 	Content string
// 	Sig     []byte
// }

type Message uint8

const (
    MessageEvent Message = iota + 1
    MessageSubscribe
    MessageClose
)

func (s Message) String() string {

    switch s {
    case MessageEvent:
        return "EVENT"
    case MessageSubscribe:
        return "REQ"
    case MessageClose:
        return "CLOSE"
    }
	panic(fmt.Sprintf("[core] unknown message type %d", s))
}

func MessageFromBytes(b []byte) (Message, error) {

    s := strings.Trim(string(b), "\"")

    switch s {
    case "EVENT":
        return MessageEvent, nil
    case "REQ":
        return MessageSubscribe, nil
    case "CLOSE":
        return MessageClose, nil
    }
	return 0, fmt.Errorf("[core] %q is not supported message type", s)
}

func MessageFromString(s string) (Message, error) {
    switch s {
    case "EVENT":
        return MessageEvent, nil
    case "REQ":
        return MessageSubscribe, nil
    case "CLOSE":
        return MessageClose, nil
    }
	return 0, fmt.Errorf("[core] %q is not supported message type", s)
}

type Tag []string

type Event struct {
	Id      string `json:"id"`
	// Pubkey  string `json:"pubkey"`
	// CreatedAt string `json:"created_at"`
	// Kind    uint8 `json:"kind"`
	// Tags    []Tag `json:"tags"`
	// Content string `json:"content"`
	// Sig     string `json:"sig"`
}
