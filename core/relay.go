package core

type RelayService interface {
	Publish(e *Event) error
	Subscribe(id SubId, filters []*Filter) error
	Close(id SubId) error
}
