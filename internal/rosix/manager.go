package rosix

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/uos-projects/uos-kernel/kernel"
)

type handle struct {
	id kernel.ResourceID
}

type watcher struct {
	ch     chan kernel.StateEvent
	filter func(kernel.StateEvent) bool
}

type Manager struct {
	twins    kernel.TwinStore
	topo     kernel.Graph
	mu       sync.RWMutex
	handles  map[kernel.FD]*handle
	nextFD   int32
	watchers map[kernel.ResourceID][]*watcher
}

func NewManager(twins kernel.TwinStore, topo kernel.Graph) *Manager {
	return &Manager{
		twins:    twins,
		topo:     topo,
		handles:  make(map[kernel.FD]*handle),
		watchers: make(map[kernel.ResourceID][]*watcher),
	}
}

func (m *Manager) Open(_ context.Context, id kernel.ResourceID) (kernel.FD, error) {
	if _, ok := m.twins.Get(id); !ok {
		return kernel.InvalidFD, fmt.Errorf("resource %s not found", id)
	}
	fd := kernel.FD(atomic.AddInt32(&m.nextFD, 1) - 1)
	m.mu.Lock()
	m.handles[fd] = &handle{id: id}
	m.mu.Unlock()
	return fd, nil
}

func (m *Manager) Close(fd kernel.FD) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.handles[fd]; !ok {
		return fmt.Errorf("invalid fd %d", fd)
	}
	delete(m.handles, fd)
	return nil
}

func (m *Manager) resolve(fd kernel.FD) (kernel.ResourceID, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h, ok := m.handles[fd]
	if !ok {
		return "", fmt.Errorf("invalid fd %d", fd)
	}
	return h.id, nil
}

func (m *Manager) Read(_ context.Context, fd kernel.FD) (*kernel.DigitalTwin, error) {
	id, err := m.resolve(fd)
	if err != nil {
		return nil, err
	}
	twin, ok := m.twins.Get(id)
	if !ok {
		return nil, fmt.Errorf("resource %s disappeared", id)
	}
	return twin, nil
}

func (m *Manager) Write(_ context.Context, fd kernel.FD, field, value, cause, actor string) error {
	id, err := m.resolve(fd)
	if err != nil {
		return err
	}
	twin, ok := m.twins.Get(id)
	if !ok {
		return fmt.Errorf("resource %s disappeared", id)
	}
	oldValue := twin.Current[field]
	event := kernel.StateEvent{
		Timestamp: time.Now(),
		Field:     field,
		OldValue:  oldValue,
		NewValue:  value,
		Cause:     cause,
		Actor:     actor,
	}
	if err := m.twins.AppendEvent(id, event); err != nil {
		return err
	}
	m.notify(id, event)
	return nil
}

func (m *Manager) RCtl(_ context.Context, fd kernel.FD, cmd string, args map[string]string) (any, error) {
	id, err := m.resolve(fd)
	if err != nil {
		return nil, err
	}
	switch cmd {
	case "get_kind":
		twin, ok := m.twins.Get(id)
		if !ok {
			return nil, fmt.Errorf("resource %s not found", id)
		}
		return twin.Resource.Kind, nil
	case "get_attributes":
		twin, ok := m.twins.Get(id)
		if !ok {
			return nil, fmt.Errorf("resource %s not found", id)
		}
		return twin.Resource.Attributes, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}
}

func (m *Manager) History(_ context.Context, fd kernel.FD, since, until time.Time) ([]kernel.StateEvent, error) {
	id, err := m.resolve(fd)
	if err != nil {
		return nil, err
	}
	return m.twins.History(id, since, until), nil
}

func (m *Manager) Watch(_ context.Context, fd kernel.FD, filter func(kernel.StateEvent) bool) (<-chan kernel.StateEvent, error) {
	id, err := m.resolve(fd)
	if err != nil {
		return nil, err
	}
	ch := make(chan kernel.StateEvent, 64)
	w := &watcher{ch: ch, filter: filter}
	m.mu.Lock()
	m.watchers[id] = append(m.watchers[id], w)
	m.mu.Unlock()
	return ch, nil
}

func (m *Manager) notify(id kernel.ResourceID, event kernel.StateEvent) {
	m.mu.RLock()
	ws := m.watchers[id]
	m.mu.RUnlock()
	for _, w := range ws {
		if w.filter == nil || w.filter(event) {
			select {
			case w.ch <- event:
			default:
			}
		}
	}
}

func (m *Manager) Relate(_ context.Context, fd1, fd2 kernel.FD, rel kernel.RelationType) error {
	id1, err := m.resolve(fd1)
	if err != nil {
		return err
	}
	id2, err := m.resolve(fd2)
	if err != nil {
		return err
	}
	m.topo.AddEdge(kernel.Edge{From: id1, To: id2, Relation: rel})
	return nil
}

func (m *Manager) Traverse(_ context.Context, fd kernel.FD, dir kernel.Direction, rel kernel.RelationType) ([]kernel.FD, error) {
	id, err := m.resolve(fd)
	if err != nil {
		return nil, err
	}
	neighbors := m.topo.Neighbors(id, dir, rel)
	fds := make([]kernel.FD, 0, len(neighbors))
	for _, nid := range neighbors {
		nfd := kernel.FD(atomic.AddInt32(&m.nextFD, 1) - 1)
		m.mu.Lock()
		m.handles[nfd] = &handle{id: nid}
		m.mu.Unlock()
		fds = append(fds, nfd)
	}
	return fds, nil
}
