package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

// type major struct {
// 	Name    string `json:"name"`
// 	Country string `json:"country"`
// }

// album represents data about a record album.
type album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
	// Major  major   `json:"major"`
}

// albums slice to seed record album data.
var albums = []album{
	// {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99, Major: major{Name: "Sony", Country: "Japan"}},
	// {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99, Major: major{Name: "Universal", Country: "Nederlands"}},
	{ID: 0, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: 1, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: 2, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.DELETE("/albums/:id", deleteAlbumByID)

	router.Run("localhost:8083")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	newAlbum.ID = len(albums)

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	fmt.Println(id, err, reflect.TypeOf(id))

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func deleteAlbumByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	fmt.Println(id, err, reflect.TypeOf(id))

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for i, a := range albums {
		if a.ID == id {
			albums = RemoveIndex(albums, i)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "album succesfully deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func RemoveIndex(s []album, index int) []album {
	return append(s[:index], s[index+1:]...)
}
