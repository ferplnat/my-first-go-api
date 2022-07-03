package apifunctions

import (
	"fmt"
	"log"
	"my-first-go-api/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var dbConn = database.ConnectSql()

// getAlbums responds with the list of all albums as JSON.
func GetAlbums(c *gin.Context) {
	// albums slice to seed record album data.
	var albums []Album
	columns := []string{"id", "title", "artist", "price"}

	rows, err := database.SelectSql(columns, "albums.album_info", dbConn)

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
func PostAlbums(c *gin.Context) {
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
