package main

import (
	"database/sql"
	"fmt"
	"strings"

	ui "github.com/gizak/termui/v3"
)

type ProcessRecord struct {
	id, time               int
	user, host, command    string
	dbName, state, sqlText sql.NullString
}

type ProcessListController struct {
	grid *ProcessListGrid
}

func newProcessListController() UIController {
	ret := &ProcessListController{
		grid: newProcessListGrid(0, 0),
	}
	return ret
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
	ds := getDataSource()
	q := fmt.Sprintf("select ID, USER, HOST, DB, COMMAND, TIME, STATE, info from PROCESSLIST where command != 'Sleep' order by TIME desc limit %d", *count)
	rows, err := ds.Query(q)
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

	info := "sqltop version 0.1"
	info += "\nProcesses: %d total, running: %d,  using DB: %d\n"
	text := fmt.Sprintf(info, totalProcesses, totalProcesses, len(usingDBs))
	text += fmt.Sprintf("\n\nTop %d order by time desc:\n", *count)
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
