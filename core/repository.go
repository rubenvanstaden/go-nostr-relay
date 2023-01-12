package core

type EventRepository interface {
	Store(e *Event) error
	FindById(id EventId) (*Event, error)
}

type FilterRepository interface {
	// Store(filter *Filter) error
	// FindById(id EventId) (*Event, error)
}
