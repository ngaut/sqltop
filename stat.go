package main

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	stat *sync.Map
	once sync.Once
)

var (
	kProcessListSQL = "SELECT ID, USER, HOST, DB, COMMAND, TIME, STATE, INFO FROM PROCESSLIST where COMMAND != 'Sleep' ORDER BY TIME DESC LIMIT %d"
)

// stat keys
const (
	TOTAL_PROCESSES = "total_processes"
	USING_DBS       = "using_dbs"
	TOTAL_READ      = "total_read"
	TOTAL_WRITE     = "total_write"
	PROCESS_LIST    = "process_list"
)

func Stat() *sync.Map {
	once.Do(func() {
		stat = &sync.Map{}
	})
	return stat
}

type ProcessRecord struct {
	id, time               int
	user, host, command    string
	dbName, state, sqlText sql.NullString
}

func refresh() {
	// update process info
	if err := refreshProcessList(); err != nil {
		cleanExit(err)
	}
}

func refreshProcessList() error {
	q := fmt.Sprintf(kProcessListSQL, Config().NumProcessToShow)
	rows, err := DB().Query(q)
	if err != nil {
		return err
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
		return err
	}

	Stat().Store(TOTAL_PROCESSES, totalProcesses)
	Stat().Store(USING_DBS, len(usingDBs))
	Stat().Store(PROCESS_LIST, records)

	return nil
}

func refreshIOStat() error {
	return nil
}

func refreshWorker() {
	for {
		refresh()
		time.Sleep(1 * time.Second)
	}
}
