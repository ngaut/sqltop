package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	ui "gopkg.in/gizak/termui.v1"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const version = "0.1"

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	refreshUI()
}

func cleanExit() {
	ui.Close()
	exec.Command("clear").Run()
	os.Exit(0)
}

type record struct {
	id, mem, time, state        int
	user, host, dbName, command string
	sqlText                     interface{}
}

func fetchProcessInfo() string {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:4000)/INFORMATION_SCHEMA")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	q := fmt.Sprintf("select ID, USER, HOST, DB, COMMAND, TIME, STATE, MEM, info  from PROCESSLIST")
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	totalProcesses := 0
	totalMem := 0
	usingDBs := make(map[string]struct{})

	var records []record
	for rows.Next() {
		var r record
		err := rows.Scan(&r.id, &r.user, &r.host, &r.dbName, &r.command, &r.time, &r.state, &r.mem, &r.sqlText)
		if err != nil {
			log.Fatal(err)
		}
		usingDBs[strings.ToLower(r.dbName)] = struct{}{}
		records = append(records, r)
		totalProcesses++
		totalMem += r.mem
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	info := "sqltop version 0.1"
	info += "\nProcesses: %d total, running: %d  Memory: %d,  using DB: %d\n"
	text := fmt.Sprintf(info, totalProcesses, totalProcesses, totalMem, len(usingDBs))
	text += fmt.Sprintf("\n\ndetails\n")

	text += fmt.Sprintf("ID      USER         HOST            DB                COMMAND   TIME     STATE  MEM      SQL\n")

	var sb strings.Builder
	for _, r := range records {
		var sqlText string
		if r.sqlText != nil {
			sqlText = fmt.Sprintf("%s", r.sqlText)
			if len(sqlText) > 128 {
				sqlText = sqlText[:128]
			}
		}
		_, _ = fmt.Fprintf(&sb, "%-6d  %-10s  %-12s  %-20s  %-8s  %-6d  %-6d  %-6d %-10s\n",
			r.id, r.user, r.host, r.dbName, r.command, r.time, r.state, r.mem, sqlText)
	}

	return text + sb.String()
}

// refreshUI periodically refreshes the screen.
func refreshUI() {
	par := ui.NewPar("")
	par.HasBorder = false
	par.Height = ui.TermHeight()
	par.Width = ui.TermWidth()

	topViewGrid := ui.NewGrid(ui.NewRow(ui.NewCol(ui.TermWidth(), 0, par)))

	// Start with the topviewGrid by default
	ui.Body.Rows = topViewGrid.Rows
	ui.Body.Align()

	redraw := make(chan struct{})

	go func() {
		for {
			par.Text = fetchProcessInfo()

			redraw <- struct{}{}
			// update every 2 seconds
			time.Sleep(2 * time.Second)
		}
	}()

	evt := ui.EventCh()
	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && (e.Ch == 'q' || e.Key == ui.KeyCtrlC) {
				cleanExit()
			}

		case <-redraw:
			ui.Render(ui.Body)
		}
	}
}
