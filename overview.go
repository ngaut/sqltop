package main

import (
	"fmt"
	"strings"
	"sync"

	ui "github.com/gizak/termui/v3"
)

var (
	overview *sync.Map
	once     sync.Once
)

const (
	TOTAL_PROCESSES = "total_processes"
	USING_DBS       = "using_dbs"
)

func Overview() *sync.Map {
	once.Do(func() {
		overview = &sync.Map{}
	})
	return overview
}

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
	if totalProcess, ok := Overview().Load(TOTAL_PROCESSES); ok {
		fmt.Fprintf(&sb, "Processes: %d total, running %d ", totalProcess, totalProcess)
	}
	if usingDBs, ok := Overview().Load(USING_DBS); ok {
		fmt.Fprintf(&sb, "using DB: %d ", usingDBs)
	}
	sb.WriteString("\n")
	c.grid.SetText(sb.String())
}
