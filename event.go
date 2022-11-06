package kernel

import "time"

type Event interface {
	Version() int
	At() time.Time
	By() string
}

type Publisher interface {
	Publish(Event) error
}

type EventModel struct {
	EventVersion int       `json:"event_version"`
	EventAt      time.Time `json:"event_at"`
	EventBy      string    `json:"event_by"`
}

func (m *EventModel) Version() int {
	return m.EventVersion
}

func (m *EventModel) At() time.Time {
	return m.EventAt
}

func (m *EventModel) By() string {
	return m.EventBy
}
