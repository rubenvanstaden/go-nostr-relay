package core

import (
	"encoding/json"
)

type RelayService interface {
	Publish(e *Event) error
	Subscribe(id SubId, filters []*Filter) error
	Close(id SubId) error
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

func EncodeRelayMessage() {

}
