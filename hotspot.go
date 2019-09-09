package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type HotSpotGrids struct {
	grids []*ui.Grid
}

func newHotSpotGrids() *HotSpotGrids {
	termWidth, _ := ui.TerminalDimensions()
	ret := &HotSpotGrids{}
	offset := 0
	// show top10 hot regions, which means we need 5 rows with 2 columns.
	for i := 0; i < 5; i++ {
		barGrid := ui.NewGrid()
		barGrid.SetRect(0, offset, termWidth, offset+3)

		g0 := widgets.NewGauge()
		g0.Title = "Table.Test.Rec1 - Rec10"
		g0.Percent = 75
		g0.BarColor = ui.ColorYellow
		g0.BorderStyle.Fg = ui.ColorWhite
		g0.TitleStyle.Fg = ui.ColorWhite
		// 2 columns
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


type KeyRange struct {
	StartKey []byte
	EndKey   []byte
}

type HotSpotInfo struct {
	TableName string
	IdxName   string
	KeyRange  KeyRange
}

type HotSpotsController struct {
	grids *HotSpotGrids
	data  []*HotSpotInfo
}

func newHotSpotsController() UIController {
	return &HotSpotsController{
		grids: newHotSpotGrids(),
	}
}

func (c *HotSpotsController) Render() {
	c.grids.Render()
}

func (c *HotSpotsController) OnResize(payload ui.Resize) {
	c.grids.OnResize(payload)
}


func (c *HotSpotsController) UpdateData() {
	// TODO
}
