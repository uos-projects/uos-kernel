package state

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemorySnapshotStore is an in-memory implementation of SnapshotStore.
// It stores snapshot history in a simple slice.
type MemorySnapshotStore struct {
	// Storage: resourceID -> list of snapshots
	history map[string][]SnapshotRecord
	mu      sync.RWMutex
}

type SnapshotRecord struct {
	Timestamp int64
	State     interface{}
}

func NewMemorySnapshotStore() *MemorySnapshotStore {
	return &MemorySnapshotStore{
		history: make(map[string][]SnapshotRecord),
	}
}

func (m *MemorySnapshotStore) Save(ctx context.Context, resourceID string, state interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	record := SnapshotRecord{
		Timestamp: time.Now().UnixMilli(),
		State:     state, // Note: In a real backend this should be a deep copy or serialized bytes
	}

	m.history[resourceID] = append(m.history[resourceID], record)
	return nil
}

func (m *MemorySnapshotStore) Load(ctx context.Context, resourceID string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	records, ok := m.history[resourceID]
	if !ok || len(records) == 0 {
		return nil, nil // No state found
	}

	// Return the latest
	return records[len(records)-1].State, nil
}

func (m *MemorySnapshotStore) LoadAt(ctx context.Context, resourceID string, timestamp int64) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	records, ok := m.history[resourceID]
	if !ok || len(records) == 0 {
		return nil, nil
	}

	// Find the snapshot closest to but not after the timestamp
	// Simple linear search for now
	var found *SnapshotRecord
	for i := range records {
		if records[i].Timestamp <= timestamp {
			found = &records[i]
		} else {
			break
		}
	}

	if found != nil {
		return found.State, nil
	}
	return nil, fmt.Errorf("no state found before timestamp %d", timestamp)
}

func (m *MemorySnapshotStore) Snapshot() (interface{}, error) {
    return nil, nil // Not used for this interface
}

func (m *MemorySnapshotStore) Restore(snapshot interface{}) error {
    return nil // Not used yet
}
