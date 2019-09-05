package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type UIController interface {
	Render()
	OnResize(ui.Resize)
	UpdateData()
}

type ProcessListGrid struct {
	grid *ui.Grid
	par  *widgets.Paragraph
}

func newProcessListGrid() *ProcessListGrid {
	termWidth, termHeight := ui.TerminalDimensions()

	par := widgets.NewParagraph()
	par.Border = false

	grid := ui.NewGrid()
	grid.SetRect(0, 15, termWidth, termHeight)
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
	pg.grid.SetRect(0, 15, payload.Width, payload.Height-15)
	ui.Render(pg.grid)
}

func (pg *ProcessListGrid) Render() {
	ui.Render(pg.grid)
}

type HotSpotGrids struct {
	grids []*ui.Grid
}

func newHotSpotGrids() *HotSpotGrids {
	termWidth, _ := ui.TerminalDimensions()
	ret := &HotSpotGrids{}
	offset := 0
	for i := 0; i < 5; i++ {
		barGrid := ui.NewGrid()
		barGrid.SetRect(0, offset, termWidth, offset+3)

		g0 := widgets.NewGauge()
		g0.Title = "Table.Test.Rec1 - Rec10"
		g0.Percent = 75
		g0.BarColor = ui.ColorYellow
		g0.BorderStyle.Fg = ui.ColorWhite
		g0.TitleStyle.Fg = ui.ColorWhite

		barGrid.Set(
			ui.NewRow(1.0,
				ui.NewCol(0.5, g0),
				ui.NewCol(0.5, g0),
			),
		)
		ret.grids = append(ret.grids, barGrid)
		offset += 3
	}
	return ret
}

func (gs *HotSpotGrids) OnResize(payload ui.Resize) {
	offset := 0
	for _, grid := range gs.grids {
		grid.SetRect(0, offset, payload.Width, offset+3)
		ui.Render(grid)
		offset += 3
	}
}

func (gs *HotSpotGrids) Render() {
	for _, grid := range gs.grids {
		ui.Render(grid)
	}
}
