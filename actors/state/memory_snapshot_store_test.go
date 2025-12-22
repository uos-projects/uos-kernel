package state

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewMemorySnapshotStore(t *testing.T) {
	store := NewMemorySnapshotStore()
	if store == nil {
		t.Fatal("NewMemorySnapshotStore() returned nil")
	}
	if store.history == nil {
		t.Fatal("history map is nil")
	}
	if len(store.history) != 0 {
		t.Fatalf("expected empty history, got %d entries", len(store.history))
	}
}

func TestMemorySnapshotStore_Save(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "test-resource-1"
	state := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	err := store.Save(ctx, resourceID, state)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	records, ok := store.history[resourceID]
	if !ok {
		t.Fatal("resource not found in history")
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	record := records[0]
	if record.Timestamp == 0 {
		t.Fatal("timestamp is zero")
	}
	if record.State == nil {
		t.Fatal("state is nil")
	}
}

func TestMemorySnapshotStore_Load_Latest(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "test-resource-1"

	// Save first snapshot
	state1 := map[string]interface{}{"version": 1, "data": "first"}
	err := store.Save(ctx, resourceID, state1)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Wait a bit to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	// Save second snapshot
	state2 := map[string]interface{}{"version": 2, "data": "second"}
	err = store.Save(ctx, resourceID, state2)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Load should return the latest snapshot
	loaded, err := store.Load(ctx, resourceID)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if loaded == nil {
		t.Fatal("Load() returned nil")
	}

	loadedMap, ok := loaded.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", loaded)
	}

	if loadedMap["version"] != 2 {
		t.Fatalf("expected version 2, got %v", loadedMap["version"])
	}
	if loadedMap["data"] != "second" {
		t.Fatalf("expected 'second', got %v", loadedMap["data"])
	}
}

func TestMemorySnapshotStore_Load_NotFound(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "non-existent-resource"

	loaded, err := store.Load(ctx, resourceID)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if loaded != nil {
		t.Fatalf("expected nil for non-existent resource, got %v", loaded)
	}
}

func TestMemorySnapshotStore_LoadAt(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "test-resource-1"

	// Save snapshot 1
	state1 := map[string]interface{}{"version": 1, "data": "first"}
	timestamp1 := time.Now().UnixMilli()
	err := store.Save(ctx, resourceID, state1)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Wait to ensure different timestamps
	time.Sleep(20 * time.Millisecond)

	// Save snapshot 2
	state2 := map[string]interface{}{"version": 2, "data": "second"}
	timestamp2 := time.Now().UnixMilli()
	err = store.Save(ctx, resourceID, state2)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Wait again
	time.Sleep(20 * time.Millisecond)

	// Save snapshot 3
	state3 := map[string]interface{}{"version": 3, "data": "third"}
	timestamp3 := time.Now().UnixMilli()
	err = store.Save(ctx, resourceID, state3)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Test LoadAt: query snapshot 1
	loaded1, err := store.LoadAt(ctx, resourceID, timestamp1+5)
	if err != nil {
		t.Fatalf("LoadAt() failed: %v", err)
	}
	if loaded1 == nil {
		t.Fatal("LoadAt() returned nil for snapshot 1")
	}
	loaded1Map := loaded1.(map[string]interface{})
	if loaded1Map["version"] != 1 {
		t.Fatalf("expected version 1, got %v", loaded1Map["version"])
	}

	// Test LoadAt: query snapshot 2
	loaded2, err := store.LoadAt(ctx, resourceID, timestamp2+5)
	if err != nil {
		t.Fatalf("LoadAt() failed: %v", err)
	}
	if loaded2 == nil {
		t.Fatal("LoadAt() returned nil for snapshot 2")
	}
	loaded2Map := loaded2.(map[string]interface{})
	if loaded2Map["version"] != 2 {
		t.Fatalf("expected version 2, got %v", loaded2Map["version"])
	}

	// Test LoadAt: query snapshot 3 (latest)
	loaded3, err := store.LoadAt(ctx, resourceID, timestamp3+5)
	if err != nil {
		t.Fatalf("LoadAt() failed: %v", err)
	}
	if loaded3 == nil {
		t.Fatal("LoadAt() returned nil for snapshot 3")
	}
	loaded3Map := loaded3.(map[string]interface{})
	if loaded3Map["version"] != 3 {
		t.Fatalf("expected version 3, got %v", loaded3Map["version"])
	}

	// Test LoadAt: query before first snapshot (should return error)
	_, err = store.LoadAt(ctx, resourceID, timestamp1-1000)
	if err == nil {
		t.Fatal("expected error for timestamp before first snapshot")
	}
}

func TestMemorySnapshotStore_MultipleResources(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()

	// Save snapshots for resource 1
	resource1 := "resource-1"
	state1 := map[string]interface{}{"resource": "1", "value": 100}
	err := store.Save(ctx, resource1, state1)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Save snapshots for resource 2
	resource2 := "resource-2"
	state2 := map[string]interface{}{"resource": "2", "value": 200}
	err = store.Save(ctx, resource2, state2)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Load resource 1
	loaded1, err := store.Load(ctx, resource1)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	loaded1Map := loaded1.(map[string]interface{})
	if loaded1Map["resource"] != "1" {
		t.Fatalf("expected resource '1', got %v", loaded1Map["resource"])
	}

	// Load resource 2
	loaded2, err := store.Load(ctx, resource2)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	loaded2Map := loaded2.(map[string]interface{})
	if loaded2Map["resource"] != "2" {
		t.Fatalf("expected resource '2', got %v", loaded2Map["resource"])
	}
}

func TestMemorySnapshotStore_ConcurrentWrites(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "concurrent-resource"
	numGoroutines := 10
	numWritesPerGoroutine := 5

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numWritesPerGoroutine; j++ {
				state := map[string]interface{}{
					"goroutine": goroutineID,
					"write":     j,
				}
				err := store.Save(ctx, resourceID, state)
				if err != nil {
					t.Errorf("Save() failed in goroutine %d: %v", goroutineID, err)
				}
				time.Sleep(time.Millisecond) // Small delay to ensure different timestamps
			}
		}(i)
	}

	wg.Wait()

	// Verify all writes were recorded
	records, ok := store.history[resourceID]
	if !ok {
		t.Fatal("resource not found in history")
	}

	expectedCount := numGoroutines * numWritesPerGoroutine
	if len(records) != expectedCount {
		t.Fatalf("expected %d records, got %d", expectedCount, len(records))
	}

	// Verify latest snapshot exists
	latest, err := store.Load(ctx, resourceID)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if latest == nil {
		t.Fatal("Load() returned nil")
	}
}

func TestMemorySnapshotStore_ConcurrentReads(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "concurrent-read-resource"

	// Pre-populate with some snapshots
	for i := 0; i < 10; i++ {
		state := map[string]interface{}{"version": i}
		err := store.Save(ctx, resourceID, state)
		if err != nil {
			t.Fatalf("Save() failed: %v", err)
		}
		time.Sleep(time.Millisecond)
	}

	numGoroutines := 20
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			loaded, err := store.Load(ctx, resourceID)
			if err != nil {
				t.Errorf("Load() failed: %v", err)
				return
			}
			if loaded == nil {
				t.Error("Load() returned nil")
				return
			}
			loadedMap := loaded.(map[string]interface{})
			if loadedMap["version"] != 9 {
				t.Errorf("expected version 9, got %v", loadedMap["version"])
			}
		}()
	}

	wg.Wait()
}

func TestMemorySnapshotStore_LoadAt_ExactTimestamp(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "exact-timestamp-resource"

	// Save snapshot with known timestamp
	state1 := map[string]interface{}{"version": 1}
	timestamp1 := time.Now().UnixMilli()
	err := store.Save(ctx, resourceID, state1)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Wait and save another snapshot
	time.Sleep(20 * time.Millisecond)
	state2 := map[string]interface{}{"version": 2}
	err = store.Save(ctx, resourceID, state2)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Query at exact timestamp of snapshot 1
	loaded, err := store.LoadAt(ctx, resourceID, timestamp1)
	if err != nil {
		t.Fatalf("LoadAt() failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadAt() returned nil")
	}
	loadedMap := loaded.(map[string]interface{})
	if loadedMap["version"] != 1 {
		t.Fatalf("expected version 1, got %v", loadedMap["version"])
	}
}

func TestMemorySnapshotStore_LoadAt_BetweenSnapshots(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "between-snapshots-resource"

	// Save snapshot 1
	state1 := map[string]interface{}{"version": 1}
	timestamp1 := time.Now().UnixMilli()
	err := store.Save(ctx, resourceID, state1)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Wait
	time.Sleep(20 * time.Millisecond)

	// Save snapshot 2
	state2 := map[string]interface{}{"version": 2}
	timestamp2 := time.Now().UnixMilli()
	err = store.Save(ctx, resourceID, state2)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Query at timestamp between snapshot 1 and 2
	betweenTimestamp := timestamp1 + (timestamp2-timestamp1)/2
	loaded, err := store.LoadAt(ctx, resourceID, betweenTimestamp)
	if err != nil {
		t.Fatalf("LoadAt() failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadAt() returned nil")
	}
	loadedMap := loaded.(map[string]interface{})
	// Should return snapshot 1 (closest before the timestamp)
	if loadedMap["version"] != 1 {
		t.Fatalf("expected version 1 (closest before timestamp), got %v", loadedMap["version"])
	}
}

func TestMemorySnapshotStore_StateTypes(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "state-types-resource"

	// Test with different state types
	testCases := []struct {
		name  string
		state interface{}
	}{
		{"map", map[string]interface{}{"key": "value"}},
		{"slice", []int{1, 2, 3}},
		{"string", "simple string"},
		{"int", 42},
		{"struct", struct {
			Field1 string
			Field2 int
		}{"test", 123}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := store.Save(ctx, resourceID+"-"+tc.name, tc.state)
			if err != nil {
				t.Fatalf("Save() failed: %v", err)
			}

			loaded, err := store.Load(ctx, resourceID+"-"+tc.name)
			if err != nil {
				t.Fatalf("Load() failed: %v", err)
			}

			if !reflect.DeepEqual(loaded, tc.state) {
				t.Fatalf("state mismatch: expected %v, got %v", tc.state, loaded)
			}
		})
	}
}

func TestMemorySnapshotStore_SnapshotAndRestore(t *testing.T) {
	store := NewMemorySnapshotStore()

	// Snapshot() should return nil, nil
	snapshot, err := store.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot() failed: %v", err)
	}
	if snapshot != nil {
		t.Fatalf("expected nil, got %v", snapshot)
	}

	// Restore() should return nil
	err = store.Restore(nil)
	if err != nil {
		t.Fatalf("Restore() failed: %v", err)
	}
}

func TestMemorySnapshotStore_LoadAt_EmptyHistory(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "empty-resource"

	// LoadAt should return nil, nil for empty history (consistent with Load)
	loaded, err := store.LoadAt(ctx, resourceID, time.Now().UnixMilli())
	if err != nil {
		t.Fatalf("LoadAt() failed: %v", err)
	}
	if loaded != nil {
		t.Fatalf("expected nil for empty history, got %v", loaded)
	}
}

func TestMemorySnapshotStore_MultipleSnapshots_Order(t *testing.T) {
	store := NewMemorySnapshotStore()
	ctx := context.Background()
	resourceID := "ordered-resource"

	// Save multiple snapshots
	timestamps := make([]int64, 5)
	for i := 0; i < 5; i++ {
		state := map[string]interface{}{"order": i}
		timestamps[i] = time.Now().UnixMilli()
		err := store.Save(ctx, resourceID, state)
		if err != nil {
			t.Fatalf("Save() failed: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Verify timestamps are in order
	records := store.history[resourceID]
	if len(records) != 5 {
		t.Fatalf("expected 5 records, got %d", len(records))
	}

	for i := 1; i < len(records); i++ {
		if records[i].Timestamp <= records[i-1].Timestamp {
			t.Fatalf("timestamps not in order: record %d timestamp %d <= record %d timestamp %d",
				i, records[i].Timestamp, i-1, records[i-1].Timestamp)
		}
	}

	// Verify LoadAt returns correct snapshots
	for i := 0; i < 5; i++ {
		loaded, err := store.LoadAt(ctx, resourceID, timestamps[i]+5)
		if err != nil {
			t.Fatalf("LoadAt() failed for snapshot %d: %v", i, err)
		}
		loadedMap := loaded.(map[string]interface{})
		if loadedMap["order"] != i {
			t.Fatalf("expected order %d, got %v", i, loadedMap["order"])
		}
	}
}
