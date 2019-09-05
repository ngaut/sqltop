package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/juju/errors"
)

type DataSource struct {
	db              *sql.DB
	dsn             string
	user, pwd, host string
	port            int
}

const (
	MaxRetryNum = 10
)

var (
	globalDS *DataSource
)

func newDataSource(user, pwd, host string, port int) *DataSource {
	return &DataSource{
		user: user,
		pwd:  pwd,
		host: host,
		port: port,
	}
}

func getDataSource() *DataSource {
	return globalDS
}

func (ds *DataSource) Close() {
	if ds.db != nil {
		if err := ds.db.Close(); err != nil {
			cleanExit(err)
		}
		ds.db = nil
	}
}

func (ds *DataSource) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/INFORMATION_SCHEMA", ds.user, ds.pwd, ds.host, ds.port)
	var err error
	ds.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return nil
}

// make sure call Connect() before calling Query()
func (ds *DataSource) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var err error
	var ret *sql.Rows
	for i := 0; i < MaxRetryNum; i++ {
		ret, err = ds.db.Query(query, args...)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			ds.db.Close()
			if err := ds.Connect(); err != nil {
				return nil, errors.Trace(err)
			}
			continue
		} else {
			return ret, nil
		}
	}
	return nil, err
}
