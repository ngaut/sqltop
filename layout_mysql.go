package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type mysqlLayout struct {
	rootGrid *ui.Grid

	processListWidget Widget
	overviewWidget    Widget

	processList *widgets.Paragraph
	overview    *widgets.Paragraph
}

func newMysqlLayout(overviewWidget Widget, processListWidget Widget) Layout {
	termWidth, termHeight := ui.TerminalDimensions()
	parent := ui.NewGrid()
	parent.SetRect(1, 3, termWidth-1, termHeight)
	parent.Border = true

	o := widgets.NewParagraph()
	o.Text = fmt.Sprintf("sqltop v0.2 [%s]\n", DB().Type())
	o.Border = false

	o.SetRect(1, 1, termWidth-1, 2)

	p := widgets.NewParagraph()
	p.Border = false

	parent.Set(
		ui.NewRow(1.0, p),
	)

	ret := &mysqlLayout{
		rootGrid:    parent,
		processList: p,
		overview:    o,

		processListWidget: processListWidget,
		overviewWidget:    overviewWidget,
	}

	return ret
}

func (l *mysqlLayout) Refresh() error {
	l.overview.Text = fmt.Sprintf("sqltop v0.2\t Mode: %s, %s\n", DB().Type(), l.overviewWidget.GetText())
	l.processList.Text = l.processListWidget.GetText()
	return nil
}

func (l *mysqlLayout) Render() {
	ui.Render(l.overview)
	ui.Render(l.rootGrid)
}

func (l *mysqlLayout) OnResize(payload interface{}) {
	resize := payload.(ui.Resize)
	l.rootGrid.SetRect(1, 3, resize.Width-1, resize.Height)
	l.overview.SetRect(1, 1, resize.Width-1, 2)
	l.Render()
}
