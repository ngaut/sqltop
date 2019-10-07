package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type DBType int

const (
	TypeUnknown DBType = iota
	TypeMySQL
	TypeTiDB
)

func (t DBType) String() string {
	switch t {
	case TypeMySQL:
		return "MySQL"
	case TypeTiDB:
		return "TiDB"
	default:
		return "Unknown"
	}
}

type DataSource struct {
	db              *sql.DB
	dsn             string
	user, pwd, host string
	port            int

	dbType DBType
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

func InitDB() error {
	cfg := Config()
	globalDS = newDataSource(cfg.DBUser, cfg.DBPwd, cfg.Host, cfg.Port)
	if err := globalDS.Connect(); err != nil {
		return err
	}
	err := globalDS.Ping()
	if err != nil {
		return err
	}
	t, err := globalDS.GetDBType()
	if err != nil {
		return err
	}
	globalDS.dbType = t
	return nil
}

func DB() *DataSource {
	return globalDS
}

func (ds *DataSource) Type() DBType {
	return ds.dbType
}

func (ds *DataSource) Ping() error {
	if ds.db == nil {
		err := ds.Connect()
		if err != nil {
			return err
		}
	}
	return ds.db.Ping()
}

func (ds *DataSource) GetDBType() (DBType, error) {
	r, err := ds.Query(`SHOW VARIABLES LIKE "version"`)
	if err != nil {
		return TypeUnknown, err
	}

	defer r.Close()
	r.Next()

	var varName, versionText string
	err = r.Scan(&varName, &versionText)
	if err != nil {
		return TypeUnknown, err
	}

	if strings.Contains(versionText, "TiDB") {
		return TypeTiDB, nil
	}
	return TypeMySQL, nil
}

func (ds *DataSource) Close() error {
	if ds.db != nil {
		if err := ds.db.Close(); err != nil {
			return err
		}
		ds.db = nil
	}
	return nil
}

func (ds *DataSource) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/INFORMATION_SCHEMA", ds.user, ds.pwd, ds.host, ds.port)
	var err error
	ds.db, err = sql.Open("mysql", dsn)
	ds.db.SetMaxIdleConns(10)
	ds.db.SetMaxOpenConns(10)

	if err != nil {
		return err
	}
	return nil
}

// make sure call Connect() before calling Query()
func (ds *DataSource) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var err error
	var ret *sql.Rows

	if ds.db == nil {
		err := ds.Connect()
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < MaxRetryNum; i++ {
		ret, err = ds.db.Query(query, args...)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			ds.db.Close()
			if err := ds.Connect(); err != nil {
				return nil, err
			}
		} else {
			return ret, nil
		}
	}
	// excees max retry, but still got error
	return nil, err
}
