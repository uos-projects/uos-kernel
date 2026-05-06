package importer

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/uos-projects/uos-kernel/internal/rosix"
	"github.com/uos-projects/uos-kernel/internal/topo"
	"github.com/uos-projects/uos-kernel/internal/twin"
	"github.com/uos-projects/uos-kernel/kernel"
)

const (
	colPerson       = 0
	colDefectType   = 1
	colArea         = 2
	colDevice       = 3
	colSwitchSt     = 4
	colFeeder       = 5
	colSubstation   = 6
	colCreated      = 7
	colAccepted     = 8
	colRepairStart  = 9
	colRepairEnd    = 10
	colArchived     = 11
	colStatus       = 12
	colCause        = 13
	colResolution   = 14
)

type Stats struct {
	Persons       int
	Areas         int
	Substations   int
	Feeders       int
	SwitchStations int
	Terminals     int
	Tickets       int
	Edges         int
}

func LoadDefects(csvPath string) (*kernel.WorldModel, *Stats, error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("read csv: %w", err)
	}
	if len(records) < 2 {
		return nil, nil, fmt.Errorf("csv has no data rows")
	}
	rows := records[1:]

	store := twin.NewMemoryTwinStore()
	graph := topo.NewMemoryGraph()
	stats := &Stats{}

	seen := make(map[kernel.ResourceID]bool)

	ensure := func(id kernel.ResourceID, kind kernel.ResourceKind, name string, attrs map[string]string) {
		if seen[id] {
			return
		}
		seen[id] = true
		t := &kernel.DigitalTwin{
			Resource: kernel.Resource{
				ID:         id,
				Kind:       kind,
				Name:       name,
				Attributes: attrs,
			},
			Current: make(map[string]string),
		}
		store.Put(t)
		switch kind {
		case kernel.KindPerson:
			stats.Persons++
		case kernel.KindArea:
			stats.Areas++
		case kernel.KindSubstation:
			stats.Substations++
		case kernel.KindFeeder:
			stats.Feeders++
		case kernel.KindSwitchStation:
			stats.SwitchStations++
		case kernel.KindTerminal:
			stats.Terminals++
		}
	}

	addEdge := func(from, to kernel.ResourceID, rel kernel.RelationType) {
		graph.AddEdge(kernel.Edge{From: from, To: to, Relation: rel})
		stats.Edges++
	}

	edgeSeen := make(map[string]bool)
	addEdgeOnce := func(from, to kernel.ResourceID, rel kernel.RelationType) {
		key := string(from) + "|" + string(to) + "|" + string(rel)
		if edgeSeen[key] {
			return
		}
		edgeSeen[key] = true
		addEdge(from, to, rel)
	}

	for i, row := range rows {
		if len(row) < 15 {
			continue
		}

		personName := row[colPerson]
		area := row[colArea]
		device := row[colDevice]
		switchSt := row[colSwitchSt]
		feeder := row[colFeeder]
		substation := row[colSubstation]

		personID := kernel.ResourceID("person:" + personName)
		areaID := kernel.ResourceID("area:" + area)
		deviceID := kernel.ResourceID("dev:" + device)
		ticketID := kernel.ResourceID(fmt.Sprintf("ticket:%05d", i+1))

		ensure(personID, kernel.KindPerson, personName, map[string]string{"area": area})
		ensure(areaID, kernel.KindArea, area, nil)
		if device != "" {
			ensure(deviceID, kernel.KindTerminal, device, map[string]string{
				"feeder":     feeder,
				"substation": substation,
			})
		}

		var substationID, feederID, switchStID kernel.ResourceID
		if substation != "" {
			substationID = kernel.ResourceID("substation:" + substation)
			ensure(substationID, kernel.KindSubstation, substation, nil)
		}
		if feeder != "" {
			feederID = kernel.ResourceID("feeder:" + feeder)
			ensure(feederID, kernel.KindFeeder, feeder, map[string]string{"substation": substation})
		}
		if switchSt != "" {
			switchStID = kernel.ResourceID("switchst:" + switchSt)
			ensure(switchStID, kernel.KindSwitchStation, switchSt, map[string]string{"feeder": feeder})
		}

		// Topology edges (deduplicated)
		addEdgeOnce(areaID, personID, kernel.RelContains)
		if substationID != "" {
			addEdgeOnce(areaID, substationID, kernel.RelContains)
		}
		if substationID != "" && feederID != "" {
			addEdgeOnce(substationID, feederID, kernel.RelContains)
		}
		if feederID != "" && switchStID != "" {
			addEdgeOnce(feederID, switchStID, kernel.RelContains)
		}
		if switchStID != "" && deviceID != "" {
			addEdgeOnce(switchStID, deviceID, kernel.RelContains)
		} else if feederID != "" && deviceID != "" && switchSt == "" {
			addEdgeOnce(feederID, deviceID, kernel.RelContains)
		}

		// Ticket resource + timeline
		ticketTwin := &kernel.DigitalTwin{
			Resource: kernel.Resource{
				ID:   ticketID,
				Kind: kernel.KindDefectTicket,
				Name: fmt.Sprintf("%s-%s-%s", area, row[colDefectType], device),
				Attributes: map[string]string{
					"defect_type": row[colDefectType],
					"cause":       row[colCause],
					"resolution":  row[colResolution],
					"status":      row[colStatus],
				},
			},
			Current: make(map[string]string),
		}
		store.Put(ticketTwin)
		stats.Tickets++

		appendStatus := func(ts, status, cause string) {
			t := parseTime(ts)
			if t.IsZero() {
				return
			}
			_ = store.AppendEvent(ticketID, kernel.StateEvent{
				Timestamp: t,
				Field:     "status",
				OldValue:  ticketTwin.Current["status"],
				NewValue:  status,
				Cause:     cause,
				Actor:     personName,
			})
		}

		appendStatus(row[colCreated], "created", "系统告警")
		appendStatus(row[colAccepted], "accepted", personName+"接单")
		appendStatus(row[colRepairStart], "repairing", "开始消缺")
		appendStatus(row[colRepairEnd], "repaired", row[colResolution])
		appendStatus(row[colArchived], "archived", "验收归档")

		// Ticket edges
		addEdge(ticketID, personID, kernel.RelAssigned)
		if deviceID != "" {
			addEdge(ticketID, deviceID, kernel.RelLocatedAt)
		}
	}

	mgr := rosix.NewManager(store, graph)

	world := &kernel.WorldModel{
		Twins: store,
		Topo:  graph,
		Rosix: mgr,
	}
	return world, stats, nil
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	layouts := []string{
		"2006-01-02 15:04:05.000000",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}
