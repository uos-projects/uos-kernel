package state

import "context"

// SnapshotStore defines the interface for persisting and retrieving actor snapshots.
// This abstraction allows for different backend implementations (Iceberg, File, Memory, etc.).
type SnapshotStore interface {
	// Save persists the snapshot of a resource.
	// resourceID: The unique identifier of the resource.
	// state: The state data to persist (usually a map or struct).
	Save(ctx context.Context, resourceID string, state interface{}) error

	// Load retrieves the latest state of a resource.
	// Returns nil if no state exists.
	Load(ctx context.Context, resourceID string) (interface{}, error)

	// LoadAt retrieves the state of a resource at a specific point in time.
	// timestamp: Unix timestamp in milliseconds.
	// Returns generic error if time travel is not supported by the backend.
	LoadAt(ctx context.Context, resourceID string, timestamp int64) (interface{}, error)
    
    // Snapshot creates a versioned snapshot of the current state
    Snapshot() (interface{}, error)
    
    // Restore restores state from a snapshot
    Restore(snapshot interface{}) error
}
