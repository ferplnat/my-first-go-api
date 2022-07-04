package apifunctions

import (
	"fmt"
	"my-first-go-api/database"
	"net/http"
	"reflect"
	"strconv"
	"strings"

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
	uriId := c.Query("id")

	// Convert to int/validate that the uri parameter is an int (only acceptable value for id)
	id, err := strconv.Atoi(uriId)

	if err != nil {
		panic(err)
	}

	// Call BindJSON to bind the received JSON to
	// updateAlbum
	var updateAlbum Album
	if err := c.BindJSON(&updateAlbum); err != nil {
		return
	}

	// Begin constructing SQL Query
	sqlQuery := "UPDATE albums.album_info SET "

	// Prepare to get field and value info from struct
	albumInfo := reflect.ValueOf(updateAlbum)
	albumItem := albumInfo.Type()

	// Create slice to store values that will be updated
	updateValues := []any{}

	// Variable to track how many iterations have been skipped, since we're using a struct
	// all of the fields will always be present. Thus, there will always be the same amount
	// of iterations. We want to skip blank/unpopulated values.
	skipped := 0
	for i := 0; i < albumInfo.NumField(); i++ {
		// Get field and value info from struct
		field := strings.ToLower(albumItem.Field(i).Name)
		value := fmt.Sprint(albumInfo.Field(i).Interface())

		// Check if value is blank, or if the field is 'id'. 'id' is parsed from Uri params
		if strings.TrimSpace(value) != "" && field != "id" {
			// If it's the first iteration, don't append a comma
			if i-skipped != 0 {
				sqlQuery = sqlQuery + ", "
			}
			updateValues = append(updateValues, string(value))
			// Probably a better way to do this, but creating the string for interpolation
			// to preserve the sanitization functionality of db.Exec() I am still creating
			// the interpolated values in the sqlQuery. Adding 1 so that the interpolation
			// definition is never 0, subtracting skipped iterations to make sure that the
			// interpolation remains sequential, otherwise, it won't work. Pray I am smart
			sqlQuery = sqlQuery + fmt.Sprintf("%s = $%d", field, i+1-skipped)
		} else {
			// If value is blank/not populated OR if field is 'id', consider it skipped
			skipped++
		}
	}
	sqlQuery = sqlQuery + fmt.Sprintf(" WHERE id = %d", id)

	fmt.Println(sqlQuery)

	res, err := dbConn.Exec(sqlQuery, updateValues...)

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
