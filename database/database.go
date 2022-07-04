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

func ValidateId(id int, table string, dbConn *sql.DB) bool {
	var ids []int
	sqlQuery := fmt.Sprintf("SELECT id FROM %s WHERE id = %d", table, id)
	rows, err := dbConn.Query(sqlQuery)

	if err != nil {
		return false
	}
	if rows == nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var idRes int
		if err := rows.Scan(&idRes); err != nil {
			return false
		}
		ids = append(ids, idRes)
	}
	if ids != nil {
		return true
	}

	return false
}
