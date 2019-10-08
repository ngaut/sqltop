package main

type Layout interface {
	Render()
	Refresh() error
	OnResize(payload interface{})
}
