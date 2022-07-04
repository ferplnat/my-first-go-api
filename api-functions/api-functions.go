package apifunctions

import (
	"fmt"
	"my-first-go-api/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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
	columns := []string{"*"}

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

	res, err := dbConn.Exec("INSERT INTO albums.album_info (artist, title, price) VALUES ($1, $2, $3)", newAlbum.Artist, newAlbum.Title, newAlbum.Price)

	if err != nil {
		panic(err.Error())
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, fmt.Sprintf("%d\n", rowCnt))
}

// UpdateAlbum updates records for an album, specified by "id" in Uri parameters
func UpdateAlbum(c *gin.Context) {
	// Parse uri for parameters
	uriQuery := c.Query("id")

	// Convert to int/validate that the uri parameter is an int (only acceptable value for id)
	id, err := strconv.Atoi(uriQuery)

	if err != nil {
		panic(err)
	}

	// Call BindJSON to bind the received JSON to
	// updateAlbum
	var updateAlbum Album
	if err := c.BindJSON(&updateAlbum); err != nil {
		return
	}

	res, err := dbConn.Exec("UPDATE albums.album_info SET artist = $1, title = $2, price = $3 WHERE id = $4", updateAlbum.Artist, updateAlbum.Title, updateAlbum.Price, id)

	if err != nil {
		panic(err.Error())
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rowCnt == 0 {
		c.IndentedJSON(http.StatusNotModified, "0 records modified. id may be invalid")
	} else {
		c.IndentedJSON(http.StatusAccepted, fmt.Sprintf("Updated %d record.\n", rowCnt))
	}
}
