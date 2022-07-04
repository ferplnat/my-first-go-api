package database

import (
	"database/sql"
	"fmt"
	"my-first-go-api/config"
	"strings"

	_ "github.com/lib/pq"
)

func ConnectSql() (db *sql.DB) {
	conf := config.LoadConfiguration()

	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Database.DBHost, conf.Database.DBPort, conf.Database.DBUser, conf.Database.DBPassword, conf.Database.DBName)

	db, err := sql.Open(conf.Database.DBType, sqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully connected to: %s"+"\n", conf.Database.DBHost)

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
