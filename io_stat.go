package main

import (
	"fmt"
	"strings"

	ui "github.com/gizak/termui/v3"
)

type IOStatController struct {
	grid *TextGrid
}

func newIOStatController() *IOStatController {
	return &IOStatController{
		grid: newTextGrid(0, 3, 15),
	}
}

func (c *IOStatController) Render() {
	c.grid.Render()
}

func (c *IOStatController) OnResize(payload ui.Resize) {
	c.grid.OnResize(payload)
}

func (c *IOStatController) UpdateData() {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Top hotspots\n")
	if list, ok := Stat().Load(TABLES_IO_STATUS); ok {
		for _, r := range list.([]TableRegionStatus) {
			fmt.Fprintf(&sb, "%s\n", r)
		}
	}
	sb.WriteString("\n")
	c.grid.SetText(sb.String())
}
