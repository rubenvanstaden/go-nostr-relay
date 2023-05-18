package core

type Author string

type Filter struct {
	Ids     []EventId `json:"ids"`
	Authors []Author  `json:"authors"`
}

type ClientEvent struct {
	SubId SubId  `json:"sub_id"`
	Event *Event `json:"event"`
}

type ClientRequest struct {
	SubId   SubId     `json:"sub_id"`
	Filters []*Filter `json:"filter"`
}
