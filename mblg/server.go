package mblg

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

func StartServer() {
	AddHandlers("/")

	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func AddHandlers(root string) {
	http.HandleFunc(root, homeHandler)
	http.HandleFunc(fmt.Sprintf("%s%s", root, "data/"), dataHandler)
}

//go:embed templates
var templates embed.FS

// homeHandler renders the homepage with a form
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Using embed.FS with template.ParseFS: https://www.reddit.com/r/golang/comments/1fllizl/comment/lo69j1e/
	tmpl, err := template.ParseFS(templates, "templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	var sheetRange = "mblg!A1:C100"
	var apiURL = "https://sheets.googleapis.com/v4/spreadsheets/" + os.Getenv("MBLG_SHEET_ID") + "/values/" + sheetRange + "?key=" + os.Getenv("GSHEETS_API_KEY")

	response, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch data from external API", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
