package main

import (
	"fmt"
	"strings"

	ui "github.com/gizak/termui/v3"
)

type OverviewController struct {
	grid *TextGrid
}

func newOverviewController() *OverviewController {
	return &OverviewController{
		grid: newTextGrid(0, 0, 4),
	}
}

func (c *OverviewController) Render() {
	c.grid.Render()
}

func (c *OverviewController) OnResize(payload ui.Resize) {
	c.grid.OnResize(payload)
}

func (c *OverviewController) UpdateData() {
	var sb strings.Builder
	fmt.Fprintf(&sb, "sqltop version 0.1\n")
	if totalProcess, ok := Stat().Load(TOTAL_PROCESSES); ok {
		fmt.Fprintf(&sb, "Processes: %d total, running %d ", totalProcess, totalProcess)
	}
	if usingDBs, ok := Stat().Load(USING_DBS); ok {
		fmt.Fprintf(&sb, "using DB: %d ", usingDBs)
	}
	sb.WriteString("\n")
	c.grid.SetText(sb.String())
}
