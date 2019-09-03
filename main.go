package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const version = "0.1"

var (
	host  = flag.String("h", "127.0.0.1", "host")
	pwd   = flag.String("p", "", "pwd")
	user  = flag.String("u", "root", "user")
	port  = flag.Int("P", 3306, "port")
	count = flag.Int("n", 50, "Number of process to show")
)

func main() {
	flag.Parse()
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	refreshUI()
}

func cleanExit(err error) {
	ui.Close()
	exec.Command("clear").Run()
	if err != nil {
		log.Print(err)
	}
	os.Exit(0)
}

type record struct {
	id, time               int
	user, host, command    string
	dbName, state, sqlText sql.NullString
}

func fetchProcessInfo() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/INFORMATION_SCHEMA", *user, *pwd, *host, *port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		cleanExit(err)
	}
	defer db.Close()
	q := fmt.Sprintf("select ID, USER, HOST, DB, COMMAND, TIME, STATE, info from PROCESSLIST where command != 'Sleep' order by TIME desc limit %d", *count)
	rows, err := db.Query(q)
	if err != nil {
		cleanExit(err)
	}
	defer rows.Close()

	totalProcesses := 0
	usingDBs := make(map[string]struct{})

	var records []record
	for rows.Next() {
		var r record
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

// refreshUI periodically refreshes the screen.
func refreshUI() {
	pg := newProcessListGrid()

	redraw := make(chan struct{})
	go func() {
		for {
			pg.SetText(fetchProcessInfo())

			redraw <- struct{}{}
			// update every 2 seconds
			time.Sleep(2 * time.Second)
		}
	}()

	evt := ui.PollEvents()
	for {
		select {
		case e := <-evt:
			if e.Type == ui.KeyboardEvent && (e.ID == "q" || e.ID == "<C-c>") {
				cleanExit(nil)
			}
			if e.ID == "<Resize>" {
				payload := e.Payload.(ui.Resize)
				pg.OnResize(payload)
			}

		case <-redraw:
			pg.Render()
		}
	}
}
