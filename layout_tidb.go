package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type tidbLayout struct {
	rootGrid *ui.Grid

	processListWidget Widget
	overviewWidget    Widget
	iostatWidget      Widget

	overview    *widgets.Paragraph
	processList *widgets.Paragraph
	iostat      *widgets.Paragraph
}

func newTiDBLayout(overviewWidget, processListWidget, iostatWidget Widget) Layout {
	termWidth, termHeight := ui.TerminalDimensions()
	parent := ui.NewGrid()
	parent.SetRect(1, 3, termWidth-1, termHeight)
	parent.Border = true

	// overview
	o := widgets.NewParagraph()
	o.Text = fmt.Sprintf("sqltop v0.2 [%s]\n", DB().Type())
	o.Border = false
	o.SetRect(1, 1, termWidth-1, 2)

	// process list
	p := widgets.NewParagraph()
	p.Border = false

	// iostat
	i := widgets.NewParagraph()
	i.Border = false

	parent.Set(
		ui.NewRow(0.3, i),
		ui.NewRow(0.7, p),
	)

	ret := &tidbLayout{
		rootGrid:    parent,
		processList: p,
		overview:    o,
		iostat:      i,

		processListWidget: processListWidget,
		overviewWidget:    overviewWidget,
		iostatWidget:      iostatWidget,
	}

	return ret
}

func (l *tidbLayout) Refresh() error {
	l.overview.Text = fmt.Sprintf("sqltop v0.2\t Mode: %s, %s\n", DB().Type(), l.overviewWidget.GetText())
	l.processList.Text = l.processListWidget.GetText()
	l.iostat.Text = l.iostatWidget.GetText()
	return nil
}

func (l *tidbLayout) Render() {
	ui.Render(l.overview)
	ui.Render(l.rootGrid)
}

func (l *tidbLayout) OnResize(payload interface{}) {
	resize := payload.(ui.Resize)
	l.rootGrid.SetRect(1, 3, resize.Width-1, resize.Height)
	l.overview.SetRect(1, 1, resize.Width-1, 2)
	l.Render()
}
