package main

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

var (
	stat *sync.Map
	once sync.Once
)

var (
	kProcessListSQL = `SELECT 
							ID, USER, HOST, DB, COMMAND, TIME, STATE, INFO
						FROM
							PROCESSLIST 
						WHERE
							COMMAND != 'Sleep' 
						ORDER BY TIME DESC LIMIT %d`

	kIOStatSQL = `SELECT 
						table_name, index_name, SUM(written_bytes) AS w, SUM(read_bytes) AS r 
					FROM
						TIKV_REGION_STATUS 
					WHERE
						db_name != "INFORMATION_SCHEMA" AND db_name != "PERFORMANCE_SCHEMA" AND db_name != "mysql" 
					GROUP BY table_name, index_name ORDER BY w DESC LIMIT 10`
)

// stat keys
const (
	TOTAL_PROCESSES  = "total_processes"
	USING_DBS        = "using_dbs"
	TOTAL_READ       = "total_read"
	TOTAL_WRITE      = "total_write"
	TABLES_IO_STATUS = "table_io_status"
	PROCESS_LIST     = "process_list"
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

type TableRegionStatus struct {
	tableName, indexName sql.NullString
	wbytes               uint64
	rbytes               uint64
}

func (r TableRegionStatus) String() string {
	return fmt.Sprintf("Table: %-20s Index: %-20s Write: %-10s Read: %-10s",
		r.tableName.String, r.indexName.String, bytefmt.ByteSize(r.wbytes), bytefmt.ByteSize(r.rbytes))
}

func refresh() {
	// update process info
	if err := refreshProcessList(); err != nil {
		cleanExit(err)
	}

	// TiDB only-feature
	if DB().Type() == TypeTiDB {
		if err := refreshIOStat(); err != nil {
			cleanExit(err)
		}
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
			return err
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
	rows, err := DB().Query(kIOStatSQL)
	if err != nil {
		return err
	}
	defer rows.Close()
	var records []TableRegionStatus
	var totalRead, totalWrite uint64
	for rows.Next() {
		var r TableRegionStatus
		err := rows.Scan(&r.tableName, &r.indexName, &r.wbytes, &r.rbytes)
		if err != nil {
			return err
		}
		records = append(records, r)

		totalRead += r.rbytes
		totalWrite += r.wbytes
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	Stat().Store(TOTAL_READ, totalRead)
	Stat().Store(TOTAL_WRITE, totalWrite)
	Stat().Store(TABLES_IO_STATUS, records)

	return nil
}

func refreshWorker() {
	for {
		refresh()
		time.Sleep(1 * time.Second)
	}
}
