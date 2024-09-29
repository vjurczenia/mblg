package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func index(c *gin.Context) {
	// The range of data you want to fetch (e.g., "Sheet1!A1:D10")
	var sheetRange = "mblg!A1:C100"

	// Construct the API URL
	var apiURL = "https://sheets.googleapis.com/v4/spreadsheets/" + os.Getenv("MBLG_SHEET_ID") + "/values/" + sheetRange + "?key=" + os.Getenv("GSHEETS_API_KEY")

	// Call the external API
	response, err := http.Get(apiURL)
	if err != nil {
		// Handle the error if the request fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch data from external API",
		})
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read response body",
		})
		return
	}

	// Unmarshal the response body into a map
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse JSON",
		})
		return
	}

	// Return the external API response as JSON
	c.JSON(http.StatusOK, result)
}

func keepAwakeOnRender() {
	for {
		http.Get("https://mblg.onrender.com/ping")
		time.Sleep(60 * time.Second)
	}
}

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	go keepAwakeOnRender()

	r := gin.Default()
	store := persistence.NewInMemoryStore(time.Minute)
	r.Use(cors.Default())

	r.GET("/", cache.CachePage(store, 1*time.Second, index))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
