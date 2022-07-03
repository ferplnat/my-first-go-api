package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

const (
	pqhost     string = "localhost"
	pqport     int    = 49153
	pquser     string = "postgres"
	pqpassword string = "postgrespw"
	pqdbname   string = "postgres"
)

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var dbConn = connectPsql(pqhost, pqport, pquser, pqpassword, pqdbname)

func main() {
	// ensure dbConn gets closed if main function exits.
	defer dbConn.Close()
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	// albums slice to seed record album data.
	var albums []Album

	sqlQuery := `SELECT * FROM albums.album_info`
	rows, err := dbConn.Query(sqlQuery)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	if rows == nil {
		c.IndentedJSON(http.StatusBadRequest, "No rows")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist,
			&alb.Price); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
		}
		albums = append(albums, alb)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum Album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	res, err := dbConn.Exec("INSERT INTO albums.album_info VALUES(?)", newAlbum)

	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusCreated, fmt.Sprintf("ID = %d, affected = %d\n", lastId, rowCnt))
}

func connectPsql(host string, port int, user string, password string, dbname string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully connected to: %s"+"\n", host)

	return db
}
