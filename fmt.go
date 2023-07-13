package main

import (
	"fmt"
	"strings"

	"github.com/chainbound/cbctl/api"
)

func printMessageTrace(trace []*api.TraceEntry, showSource bool) {
	if showSource {
		fmt.Printf("Timestamp\tNode ID\t\t\tRegion\t\t\tObservation Type\tSource\n")
	} else {
		fmt.Printf("Timestamp\tNode ID\t\t\tRegion\t\t\tObservation Type\n")

	}
	for _, entry := range trace {
		chunks := strings.Split(entry.NodeID, "-")[:3]
		id := strings.Join(chunks, "-")
		if showSource {
			fmt.Printf("[%d]\t%s\t(%s)\t\t%s\t\t\t%s\n", entry.Timestamp/1000, id, entry.Region, entry.ObservationType, entry.Source)
		} else {
			fmt.Printf("[%d]\t%s\t(%s)\t\t%s\n", entry.Timestamp/1000, id, entry.Region, entry.ObservationType)
		}
	}
}
