package kernel

import (
	"encoding/json"
	"sync"
)

type Serializer interface {
	Bind(...Event)
	Serialize(e Event) (Record, error)
	Deserialize(rec Record) (Event, error)
}

type jsonSerializer struct {
	mu    sync.RWMutex
	types map[string]Cloner
}

func (s *jsonSerializer) Bind(events ...Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range events {
		t := eventTyper(e)

		if _, ok := s.types[t]; ok {
			continue
		}

		s.types[t] = cloner(e)
	}
}

func (s *jsonSerializer) Serialize(e Event) (Record, error) {
	r := Record{}

	r.Version = e.Version()
	r.At = e.At()
	r.By = e.By()
	r.Type = eventTyper(e)

	data, err := json.Marshal(e)

	if err != nil {
		return r, err
	}

	r.Data = data

	return r, nil
}

func (s *jsonSerializer) Deserialize(rec Record) (Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e := s.types[rec.Type].Clone().(Event)

	if err := json.Unmarshal(rec.Data, e); err != nil {
		return e, err
	}

	return e, nil
}
