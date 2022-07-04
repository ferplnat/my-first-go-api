package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

const (
	dbtype     string = "postgres"
	dbhost     string = "localhost"
	dbport     int    = 49153
	dbuser     string = "postgres"
	dbpassword string = "postgrespw"
	dbname     string = "postgres"
)

func ConnectSql() (db *sql.DB) {
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := sql.Open(dbtype, sqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully connected to: %s"+"\n", dbhost)

	return db
}

// select function, returns data rows.
func SelectSql(columns []string, table string, dbConn *sql.DB, orderBy string, desc int) (*sql.Rows, error) {
	sqlQuery := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s", strings.Join(columns[:], ","), table, orderBy)
	if desc == 1 {
		sqlQuery = sqlQuery + " DESC"
	}
	fmt.Println(sqlQuery)
	rows, err := dbConn.Query(sqlQuery)

	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}

	return rows, nil
}
