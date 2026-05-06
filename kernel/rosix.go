package kernel

import (
	"context"
	"time"
)

type FD int32

const InvalidFD FD = -1

type ROSIX interface {
	Open(ctx context.Context, id ResourceID) (FD, error)
	Close(fd FD) error
	Read(ctx context.Context, fd FD) (*DigitalTwin, error)
	Write(ctx context.Context, fd FD, field, value, cause, actor string) error
	RCtl(ctx context.Context, fd FD, cmd string, args map[string]string) (any, error)
	History(ctx context.Context, fd FD, since, until time.Time) ([]StateEvent, error)
	Watch(ctx context.Context, fd FD, filter func(StateEvent) bool) (<-chan StateEvent, error)
	Relate(ctx context.Context, fd1, fd2 FD, rel RelationType) error
	Traverse(ctx context.Context, fd FD, dir Direction, rel RelationType) ([]FD, error)
}
