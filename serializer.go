package kernel

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Serializer interface {
	Bind(...Event)
	Serialize(e Event) (Record, error)
	Deserialize(rec Record) (Event, error)
}

type jsonSerializer struct {
	mu    sync.RWMutex
	types map[string]func() Event
}

type jsonRecord struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func NewJSONSerializer() Serializer {
	return &jsonSerializer{
		mu:    sync.RWMutex{},
		types: map[string]func() Event{},
	}
}

func (s *jsonSerializer) Bind(events ...Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range events {
		t := eventTyper(e)

		if _, ok := s.types[t]; ok {
			continue
		}

		s.types[t] = eventCTOR(e)
	}
}

func (s *jsonSerializer) Serialize(e Event) (Record, error) {
	r := Record{}
	t := eventTyper(e)
	data, err := json.Marshal(e)

	if err != nil {
		return r, err
	}

	r.Data, err = json.Marshal(&jsonRecord{
		Type: t,
		Data: json.RawMessage(data),
	})

	if err != nil {
		return r, err
	}

	r.Version = e.Version()

	return r, nil
}

func (s *jsonSerializer) Deserialize(rec Record) (Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jrec := &jsonRecord{}

	if err := json.Unmarshal(rec.Data, jrec); err != nil {
		return nil, err
	}

	ctor, ok := s.types[jrec.Type]

	if !ok {
		return nil, fmt.Errorf("unknown event type: %v", jrec.Type)

	}

	e := ctor()

	if err := json.Unmarshal(jrec.Data, e); err != nil {
		return e, err
	}

	return e, nil
}
