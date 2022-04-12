package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	fmt.Println("INIT !!!!!!!!!!!!!!!")
	albumOptions := options.Client().ApplyURI("mongodb://mongo:27017/")
	album, err := mongo.Connect(ctx, albumOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = album.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = album.Database("albumsDB").Collection("albums")

	fmt.Println("DB initialized")
}

// album represents data about a record album.
type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
	// Major  major   `json:"major"`
}

// albums slice to seed record album data.
var albumsDefault = []Album{
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
	fmt.Println("Router configured")

	for _, a := range albumsDefault {
		_, err := insertAlbum(a)
		fmt.Println(err)
	}

	router.Run(":8083")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	filter := bson.D{{}}
	//collection.Find(ctx, filter)
	albums, _ := filterAlbums(filter)
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println(id, err, reflect.TypeOf(id))

	filter := bson.D{primitive.E{Key: "id", Value: id}}

	album, _ := filterAlbums(filter)

	if len(album) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, album[0])
}

func filterAlbums(filter interface{}) ([]*Album, error) {
	// A slice of tasks for storing the decoded documents
	var albums []*Album

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return albums, err
	}

	// Iterate through the cursor and decode each document one at a time
	for cur.Next(ctx) {
		var a Album
		err := cur.Decode(&a)
		if err != nil {
			return albums, err
		}

		albums = append(albums, &a)
	}

	if err := cur.Err(); err != nil {
		return albums, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	return albums, nil
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum Album

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An error happened"})
		return
	}

	newAlbum, err := insertAlbum(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An error happened"})
		return
	}

	// Add the new album to the slice.
	//albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func insertAlbum(newAlbum Album) (Album, error) {
	newAlbum.ID = countAlbums()

	_, err := collection.InsertOne(ctx, newAlbum)
	return newAlbum, err
}

func countAlbums() int64 {
	nb, _ := collection.CountDocuments(ctx, bson.D{{}})
	return nb
}

func deleteAlbumByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	fmt.Println(id, err, reflect.TypeOf(id))

	filter := bson.D{primitive.E{Key: "id", Value: id}}

	res, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An error appened"})
		return
	} else if res.DeletedCount == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "album succesfully deleted"})
	}
}

// Just for info,used before mongoisation
// func RemoveIndex(s []Album, index int) []Album {
// 	return append(s[:index], s[index+1:]...)
// }
