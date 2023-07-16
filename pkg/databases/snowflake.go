package databases

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	sf "github.com/snowflakedb/gosnowflake"
)

type SnowflakeConn struct {
	DB *sql.DB
}

func NewSnowflakeConn(config *sf.Config) (*SnowflakeConn, error) {
	dsn, err := sf.DSN(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create DSN from Config: %v, err: %v", config, err)
	}

	_conn, err := sql.Open("snowflake", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Snowflake DB: %v", err)
	}

	// // send select 1
	// err = _conn.Ping()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to ping Snowflake DB: %v", err)
	// }

	conn := &SnowflakeConn{
		DB: _conn,
	}

	return conn, nil
}

func (conn *SnowflakeConn) RunQuery(query string) (*sql.Rows, error) {
	rows, err := conn.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %v, err: %v", query, err)
	}
	return rows, nil
}

func (conn *SnowflakeConn) PreviewRowResults(rowResults []any) {
	rowResult := rowResults[0]
	values := reflect.ValueOf(rowResult)
	columns := values.Type()
	row := make(map[string]interface{})

	for i := 0; i < values.NumField(); i++ {
		column := columns.Field(i).Name
		value := values.Field(i).Interface()
		row[column] = value
	}

	maxLength := 0
	for key := range row {
		if len(key) > maxLength {
			maxLength = len(key)
		}
	}

	for key, value := range row {
		spaces := strings.Repeat(" ", maxLength-len(key))
		log.Printf("Field: %s Value: %v", key+spaces, value)
	}
}
