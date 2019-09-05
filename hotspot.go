package main

import "sync"

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
	mu    sync.RWMutex
}

func (c *HotSpotsController) Render() {
}

func (c *HotSpotsController) Update() {
}
