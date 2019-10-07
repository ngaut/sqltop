package main

import (
	"fmt"
	"strings"
)

type OverviewController struct{}

func (c *OverviewController) GenUIText() string {
	var sb strings.Builder
	if totalProcess, ok := Stat().Load(TOTAL_PROCESSES); ok {
		fmt.Fprintf(&sb, "Processes: %d total, running %d ", totalProcess, totalProcess)
	}
	if usingDBs, ok := Stat().Load(USING_DBS); ok {
		fmt.Fprintf(&sb, "using DB: %d ", usingDBs)
	}
	return sb.String()
}
