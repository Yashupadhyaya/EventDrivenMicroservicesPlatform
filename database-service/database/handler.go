package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	postgresDB *sql.DB
	mysqlDB    *sql.DB
)

func InitPostgres(dataSourceName string) error {
	var err error
	postgresDB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}
	return postgresDB.Ping()
}

func InitMySQL(dataSourceName string) error {
	var err error
	mysqlDB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	return mysqlDB.Ping()
}

func GetDB() *sql.DB {
	if postgresDB != nil {
		return postgresDB
	}
	if mysqlDB != nil {
		return mysqlDB
	}
	log.Println("no database initialized")
	return nil
}

func Query(sqlDB *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := sqlDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}

		result = append(result, row)
	}

	return result, nil
}

func Execute(sqlDB *sql.DB, query string, args ...interface{}) (int64, error) {
	result, err := sqlDB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
