package main

import (
	"fmt"
	"regexp"
	"strings"

	ui "github.com/gizak/termui/v3"
)

var (
	re = regexp.MustCompile(`\r?\n`)
)

type ProcessListController struct {
	grid *TextGrid
}

func newProcessListController() UIController {
	_, termHeight := ui.TerminalDimensions()
	return &ProcessListController{
		grid: newTextGrid(0, 20, termHeight-20),
	}
}

func (c *ProcessListController) Render() {
	c.grid.Render()
}

func (c *ProcessListController) OnResize(payload ui.Resize) {
	c.grid.OnResize(payload)
}

func (c *ProcessListController) UpdateData() {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Top %d order by time desc:\n", Config().NumProcessToShow)
	fmt.Fprintf(&sb, "%-6s  %-20s  %-20s  %-20s  %-7s  %-6s  %-8s  %-15s\n",
		"ID", "USER", "HOST", "DB", "COMMAND", "TIME", "STATE", "SQL")

	if r, ok := Stat().Load(PROCESS_LIST); ok {
		records := r.([]ProcessRecord)
		for _, r := range records {
			var sqlText string
			if r.sqlText.Valid {
				sqlText = r.sqlText.String
				if len(sqlText) > 128 {
					sqlText = sqlText[:128]
					sqlText = re.ReplaceAllString(sqlText, " ")
				}
			}
			fmt.Fprintf(&sb, "%-6d  %-20s  %-20s  %-20s  %-7s  %-6d  %-8s  %-15s\n",
				r.id, r.user, r.host, r.dbName.String, r.command, r.time, r.state.String, sqlText)
		}
	}
	c.grid.SetText(sb.String())
}
