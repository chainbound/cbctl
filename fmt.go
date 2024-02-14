package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chainbound/cbctl/api"

	"github.com/jedib0t/go-pretty/v6/table"
)

func printMessageTrace(trace []*api.TraceEntry, showSource bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)

	if showSource {
		t.AppendHeader(table.Row{"Timestamp", "Node ID", "Region", "Observation Type", "Source"})
	} else {
		t.AppendHeader(table.Row{"Timestamp", "Node ID", "Region", "Observation Type"})
	}

	firstTimestamp := trace[0].Timestamp / 1000
	fiberPropagationTime := 0
	p2pPropagationTime := 0

	for _, entry := range trace {
		id := entry.NodeID

		if entry.ObservationType == "fiber" {
			fiberPropagationTime = int(entry.Timestamp/1000) - int(firstTimestamp)
		} else if entry.ObservationType == "p2p" {
			p2pPropagationTime = int(entry.Timestamp/1000) - int(firstTimestamp)
		}

		timestamp := time.Unix(0, int64(entry.Timestamp/1000)*int64(time.Millisecond))
		timestampFmt := timestamp.Format("2006-01-02 15:04:05.000")
		if showSource {
			t.AppendRow(table.Row{timestampFmt, id, entry.Region, entry.ObservationType, entry.Source})
		} else {
			t.AppendRow(table.Row{timestampFmt, id, entry.Region, entry.ObservationType})
		}
	}

	t.AppendFooter(table.Row{"fiber propagation time", fmt.Sprint(fiberPropagationTime) + "ms"})
	t.AppendFooter(table.Row{"p2p propagation time", fmt.Sprint(p2pPropagationTime) + "ms"})

	t.Render()
}
