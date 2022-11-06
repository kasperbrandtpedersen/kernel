package kernel

import (
	"sort"
	"sync"
	"time"
)

type Record struct {
	Version int       `json:"version"`
	At      time.Time `json:"at"`
	By      string    `json:"by"`
	Type    string    `json:"type"`
	Data    []byte    `json:"data"`
}

type History []Record

// Len implements sort.Interface
func (h History) Len() int {
	return len(h)
}

// Swap implements sort.Interface
func (h History) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Less implements sort.Interface
func (h History) Less(i, j int) bool {
	return h[i].Version < h[j].Version
}

type Store interface {
	Save(stream string, recs History) error
	Load(stream string) (History, error)
}

// store implements Store as a map based storage
type store struct {
	mu      sync.RWMutex
	content map[string]History
}

func (s *store) Save(stream string, recs History) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.content[stream]; !ok {
		s.content[stream] = History{}
	}

	h := append(s.content[stream], recs...)
	sort.Sort(h)
	s.content[stream] = h

	return nil
}

func (s *store) Load(stream string) (History, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	h, ok := s.content[stream]

	if !ok {
		return History{}, nil
	}

	recs := make(History, len(h))
	recs = append(recs, h...)

	return recs, nil
}
