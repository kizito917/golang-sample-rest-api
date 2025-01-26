package main

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, albums)
}

func showServerMain(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "Welcome to my API")
}

func addNewAlbum(context *gin.Context) {
	// Create new variable to hold the type of album for new album to be added
	var newAlbum album

	if err := context.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	context.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumDetails(context *gin.Context) {
	id := context.Param("id")

	if id == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Album ID is required"})
	}

	for _, album := range albums {
		if album.ID == id {
			context.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
}

func updateAlbum(context *gin.Context) {
	id := context.Param("id")

	if id == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Album ID is required"})
	}

	var albumToUpdate map[string]interface{}

	if err := context.BindJSON(&albumToUpdate); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	for i := range albums {
		if albums[i].ID == id {
			// Get reflect values
			val := reflect.ValueOf(&albums[i]).Elem()

			// Update each field from the payload
			for key, value := range albumToUpdate {
				// Get field by name (capitalize first letter for exported fields)
				field := val.FieldByName(strings.Title(key))
				if field.IsValid() && field.CanSet() {
					// Convert and set the value
					newValue := reflect.ValueOf(value)
					if field.Type() == newValue.Type() {
						field.Set(newValue)
					}
				}
			}

			context.IndentedJSON(http.StatusOK, albums[i])
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
}

func deleteAlbum(context *gin.Context) {
	id := context.Param("id")

	if id == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Album ID is required"})
	}

	for i, album := range albums {
		if album.ID == id {
			albums = append(albums[:i], albums[i+1:]...)
			context.IndentedJSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
}

func main() {
	router := gin.Default()
	router.GET("/", showServerMain)
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumDetails)
	router.POST("/album", addNewAlbum)
	router.PATCH("/albums/:id", updateAlbum)
	router.DELETE("/albums/:id", deleteAlbum)

	router.Run("localhost:8080")
}
