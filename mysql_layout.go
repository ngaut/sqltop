package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type mysqlLayout struct {
	rootGrid *ui.Grid

	processListWidget *widgets.Paragraph
	overviewWidget    *widgets.Paragraph
}

func newMysqlLayout() Layout {
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
		rootGrid:          parent,
		processListWidget: p,
		overviewWidget:    o,
	}

	return ret
}

func (l *mysqlLayout) setOverviewText(str string) {
	l.overviewWidget.Text = fmt.Sprintf("sqltop v0.2 [%s] %s\n", DB().Type(), str)
}

func (l *mysqlLayout) setProcesslistText(str string) {
	l.processListWidget.Text = str
}

func (l *mysqlLayout) Refresh() error {
	return nil
}

func (l *mysqlLayout) Render() {
	ui.Render(l.overviewWidget)
	ui.Render(l.rootGrid)
}

func (l *mysqlLayout) OnResize(payload interface{}) {
	resize := payload.(ui.Resize)
	l.rootGrid.SetRect(1, 3, resize.Width-1, resize.Height)
	l.overviewWidget.SetRect(1, 1, resize.Width-1, 2)
	l.Render()
}
