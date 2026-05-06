package twin

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/uos-projects/uos-kernel/kernel"
)

type MemoryTwinStore struct {
	mu    sync.RWMutex
	twins map[kernel.ResourceID]*kernel.DigitalTwin
}

func NewMemoryTwinStore() *MemoryTwinStore {
	return &MemoryTwinStore{
		twins: make(map[kernel.ResourceID]*kernel.DigitalTwin),
	}
}

func (s *MemoryTwinStore) Get(id kernel.ResourceID) (*kernel.DigitalTwin, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.twins[id]
	return t, ok
}

func (s *MemoryTwinStore) Put(twin *kernel.DigitalTwin) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if twin.Current == nil {
		twin.Current = make(map[string]string)
	}
	s.twins[twin.Resource.ID] = twin
}

func (s *MemoryTwinStore) All() []*kernel.DigitalTwin {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*kernel.DigitalTwin, 0, len(s.twins))
	for _, t := range s.twins {
		result = append(result, t)
	}
	return result
}

func (s *MemoryTwinStore) AppendEvent(id kernel.ResourceID, event kernel.StateEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.twins[id]
	if !ok {
		return fmt.Errorf("resource %s not found", id)
	}
	t.Timeline = append(t.Timeline, event)
	if t.Current == nil {
		t.Current = make(map[string]string)
	}
	t.Current[event.Field] = event.NewValue
	return nil
}

func (s *MemoryTwinStore) History(id kernel.ResourceID, since, until time.Time) []kernel.StateEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.twins[id]
	if !ok {
		return nil
	}
	timeline := t.Timeline
	var start, end int
	if since.IsZero() {
		start = 0
	} else {
		start = sort.Search(len(timeline), func(i int) bool {
			return !timeline[i].Timestamp.Before(since)
		})
	}
	if until.IsZero() {
		end = len(timeline)
	} else {
		end = sort.Search(len(timeline), func(i int) bool {
			return timeline[i].Timestamp.After(until)
		})
	}
	if start >= end {
		return nil
	}
	result := make([]kernel.StateEvent, end-start)
	copy(result, timeline[start:end])
	return result
}
