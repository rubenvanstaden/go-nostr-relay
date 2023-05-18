package inmem

import (
	"fmt"
	"log"

	"github.com/rubenvanstaden/go-nostr-relay/core"
)

type repository struct {
	db map[core.EventId]*core.Event

	core.EventRepository
}

func New() core.EventRepository {
	return &repository{
		db: make(map[core.EventId]*core.Event),
	}
}

func (s *repository) Store(e *core.Event) error {

	op := "inmem.Store"

	s.db[e.Id] = e

	log.Printf("[%s] stored event: %s", op, e)

	return nil
}

func (s *repository) FindById(id core.EventId) (*core.Event, error) {

	op := "inmem.FindId"

	if e, found := s.db[id]; found {
		log.Printf("[%s] retrieved event: %s", op, e)
		return e, nil
	}

	return nil, fmt.Errorf("event not found")
}
