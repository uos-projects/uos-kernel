package kernel

import "time"

type ResourceID string

type ResourceKind string

const (
	KindSubstation    ResourceKind = "Substation"
	KindFeeder        ResourceKind = "Feeder"
	KindSwitchStation ResourceKind = "SwitchStation"
	KindTerminal      ResourceKind = "Terminal"
	KindPerson        ResourceKind = "Person"
	KindTeam          ResourceKind = "Team"
	KindArea          ResourceKind = "Area"
	KindDefectTicket  ResourceKind = "DefectTicket"
)

type Resource struct {
	ID         ResourceID
	Kind       ResourceKind
	Name       string
	Attributes map[string]string
	CreatedAt  time.Time
}
