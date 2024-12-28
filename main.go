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
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func showServerMain(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Welcome to my API")
}

func addNewAlbum(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumDetails(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Album ID is required"})
	}

	for _, album := range albums {
		if album.ID == id {
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
}

func updateAlbum(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Album ID is required"})
	}

	var albumToUpdate map[string]interface{}

	if err := c.BindJSON(&albumToUpdate); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
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

			c.IndentedJSON(http.StatusOK, albums[i])
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
}

func deleteAlbum(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Album ID is required"})
	}

	for i, album := range albums {
		if album.ID == id {
			albums = append(albums[:i], albums[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
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
