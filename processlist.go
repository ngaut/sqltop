package main

import (
	"fmt"
	"strings"

	ui "github.com/gizak/termui/v3"
)

type ProcessListController struct {
	grid *TextGrid
}

func newProcessListController() UIController {
	_, termHeight := ui.TerminalDimensions()
	return &ProcessListController{
		grid: newTextGrid(0, 3, termHeight-3),
	}
}

func (c *ProcessListController) Render() {
	c.grid.Render()
}

func (c *ProcessListController) OnResize(payload ui.Resize) {
	c.grid.OnResize(payload)
}

func (c *ProcessListController) UpdateData() {
	c.grid.SetText(c.fetchProcessInfo())
}

func (c *ProcessListController) fetchProcessInfo() string {
	q := fmt.Sprintf("select ID, USER, HOST, DB, COMMAND, TIME, STATE, info from PROCESSLIST where command != 'Sleep' order by TIME desc limit %d", *count)
	rows, err := DB().Query(q)
	if err != nil {
		cleanExit(err)
	}
	defer rows.Close()

	totalProcesses := 0
	usingDBs := make(map[string]struct{})

	var records []ProcessRecord
	for rows.Next() {
		var r ProcessRecord
		err := rows.Scan(&r.id, &r.user, &r.host, &r.dbName, &r.command, &r.time, &r.state, &r.sqlText)
		if err != nil {
			cleanExit(err)
		}
		if r.dbName.Valid {
			usingDBs[strings.ToLower(r.dbName.String)] = struct{}{}
		}
		records = append(records, r)
		totalProcesses++
	}
	err = rows.Err()
	if err != nil {
		cleanExit(err)
	}

	// update overview info
	Stat().Store(TOTAL_PROCESSES, totalProcesses)
	Stat().Store(USING_DBS, len(usingDBs))

	text := fmt.Sprintf("Top %d order by time desc:\n", *count)
	text += fmt.Sprintf("%-6s  %-20s  %-20s  %-20s  %-7s  %-6s  %-8s  %-15s\n",
		"ID", "USER", "HOST", "DB", "COMMAND", "TIME", "STATE", "SQL")

	var sb strings.Builder
	for _, r := range records {
		var sqlText string
		if r.sqlText.Valid {
			sqlText = r.sqlText.String
			if len(sqlText) > 128 {
				sqlText = sqlText[:128]
			}
		}
		_, _ = fmt.Fprintf(&sb, "%-6d  %-20s  %-20s  %-20s  %-7s  %-6d  %-8s  %-15s\n",
			r.id, r.user, r.host, r.dbName.String, r.command, r.time, r.state.String, sqlText)
	}

	return text + sb.String()
}
