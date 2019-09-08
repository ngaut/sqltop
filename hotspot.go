package main

import (

	ui "github.com/gizak/termui/v3"
)

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
