package kernel

type RelationType string

const (
	RelContains   RelationType = "contains"
	RelBelongsTo  RelationType = "belongs_to"
	RelAssigned   RelationType = "assigned"
	RelLocatedAt  RelationType = "located_at"
	RelServices   RelationType = "services"
	RelFeedsFrom  RelationType = "feeds_from"
)

type Direction int

const (
	Outbound Direction = iota
	Inbound
	Both
)

type Edge struct {
	From     ResourceID
	To       ResourceID
	Relation RelationType
	Meta     map[string]string
}

type Graph interface {
	AddEdge(edge Edge)
	RemoveEdge(from, to ResourceID, rel RelationType)
	Neighbors(id ResourceID, dir Direction, rel RelationType) []ResourceID
	Traverse(start ResourceID, dir Direction, rel RelationType, maxDepth int) []ResourceID
	EdgesOf(id ResourceID) []Edge
}
