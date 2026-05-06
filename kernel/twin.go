package kernel

import "time"

type StateEvent struct {
	Timestamp time.Time
	Field     string
	OldValue  string
	NewValue  string
	Cause     string
	Actor     string
}

type DigitalTwin struct {
	Resource Resource
	Timeline []StateEvent
	Current  map[string]string
}

type TwinStore interface {
	Get(id ResourceID) (*DigitalTwin, bool)
	Put(twin *DigitalTwin)
	All() []*DigitalTwin
	AppendEvent(id ResourceID, event StateEvent) error
	History(id ResourceID, since, until time.Time) []StateEvent
}
