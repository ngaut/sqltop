package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type ProcessListGrid struct {
	grid *ui.Grid
	par  *widgets.Paragraph
}

func newProcessListGrid() *ProcessListGrid {
	termWidth, termHeight := ui.TerminalDimensions()

	par := widgets.NewParagraph()
	par.Border = false

	grid := ui.NewGrid()
	grid.SetRect(0, 10, termWidth, termHeight-10)
	grid.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0, par),
		),
	)
	return &ProcessListGrid{
		grid: grid,
		par:  par,
	}
}

// it's caller's duty to be threaded safe
func (pg *ProcessListGrid) SetText(str string) {
	pg.par.Text = str
}

func (pg *ProcessListGrid) OnResize(payload ui.Resize) {
	pg.grid.SetRect(0, 10, payload.Width, payload.Height-10)
	ui.Render(pg.grid)
}

func (pg *ProcessListGrid) Render() {
	ui.Render(pg.grid)
}
