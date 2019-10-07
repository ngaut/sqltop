package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	re = regexp.MustCompile(`\r?\n`)
)

type ProcessListController struct {
}

func newProcessListController() *ProcessListController {
	return &ProcessListController{}
}

func (c *ProcessListController) GenUIContent() string {
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
	return sb.String()
}
