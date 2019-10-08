package main

import (
	"fmt"
	"strings"

	"code.cloudfoundry.org/bytefmt"
)

type OverviewWidget struct{}

func newOverviewWidget() Widget {
	return &OverviewWidget{}
}

func (c *OverviewWidget) GetText() string {
	var sb strings.Builder
	if totalProcess, ok := Stat().Load(TOTAL_PROCESSES); ok {
		fmt.Fprintf(&sb, "Processes: %d total, Running: %d", totalProcess, totalProcess)
	}
	if usingDBs, ok := Stat().Load(USING_DBS); ok {
		fmt.Fprintf(&sb, ", Using DB: %d", usingDBs)
	}
	if DB().Type() == TypeTiDB {
		if totalRead, ok := Stat().Load(TOTAL_READ); ok {
			fmt.Fprintf(&sb, ", Recent total read: %s", bytefmt.ByteSize(totalRead.(uint64)))
		}
		if totalWrite, ok := Stat().Load(TOTAL_WRITE); ok {
			fmt.Fprintf(&sb, ", Recent total write: %s", bytefmt.ByteSize(totalWrite.(uint64)))
		}
	}
	return sb.String()
}
