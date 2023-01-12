package relay

import (
	"log"

	"github.com/rubenvanstaden/go-nostr-relay/core"
)

type service struct {
	repository core.EventRepository

	core.RelayService
}

func New(repository core.EventRepository) core.RelayService {
	return &service{
		repository: repository,
	}
}

func (s *service) Publish(e *core.Event) error {

	// op := "relay.Publish"

	log.Printf("Publish event: %s", e)

	s.repository.Store(e)

	return nil
}

func (s *service) Subscribe(id core.SubId, filters []*core.Filter) error {
	return nil
}

func (s *service) Close(id core.SubId) error {
	return nil
}
