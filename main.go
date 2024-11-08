package main

import (
	"ascii-art-web/pkg/generator"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Title string
	Msg   string
}
type AsciiPg struct {
	P      Page
	Text   string
	Banner string
}

var indexTmpl = template.Must(template.ParseFiles("templates/index.html"))
var errTmpl = template.Must(template.ParseFiles("templates/error.html"))

func main() {
	reqFiles := []string{
		"./assets/styles/shadow.txt",
		"./assets/styles/standard.txt",
		"./assets/styles/thinkertoy.txt",
		"./static/styles.css",
	}
	checkRequired(reqFiles)

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/ascii-art", homeHandler)

	log.Println("Listening and serving on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// homeHandler() handlers GET and POST request to "/" and "/ascii-art"
func homeHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" && req.URL.Path != "/ascii-art" {
		errorHandler(w, http.StatusNotFound)
		return
	}
	page := &AsciiPg{Banner: "standard"}
	if !isFileThere("./assets/styles/" + page.Banner + ".txt") {
		errorHandler(w, http.StatusInternalServerError)
		return
	}
	page.P.Title, _ = generator.GenArt("ASCII-ART", "standard")
	page.P.Msg, _ = generator.GenArt("Enter text\nhere", "standard")

	switch req.Method {
	case "GET":
		indexTmpl.Execute(w, page)
	case "POST":
		handlePost(w, req, page)
	default:
		errorHandler(w, http.StatusMethodNotAllowed)
	}
}

// handlePost() handles
func handlePost(w http.ResponseWriter, req *http.Request, page *AsciiPg) {
	text, banner, formErr := getFormInputs(req)
	if formErr != "" {
		errorHandler(w, http.StatusBadRequest)
		return
	}
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.Trim(text, "\n")

	if !isFileThere("./assets/styles/" + banner + ".txt") {
		errorHandler(w, http.StatusInternalServerError)
		return
	}
	if art, err := generator.GenArt(text, banner); err != nil {
		page.P.Msg = "Failed to generate ASCII art: " + err.Error()
	} else {
		page.P.Msg = art
	}
	page.Text = text
	page.Banner = banner

	indexTmpl.Execute(w, page)
}

// getFormInputs() gets the text and banner input from the form in the POST request.
// Returns text, banner and error message.
func getFormInputs(req *http.Request) (string, string, string) {
	if err := req.ParseForm(); err != nil {
		return "", "", "400"
	}
	text := req.Form.Get("text")
	banner := req.Form.Get("banner")

	if text == "" || banner == "" {
		return "", "", "400"
	}
	return text, banner, ""
}

// errorHandler() generate custom error page responses.
// If error.html can't be parse, default to simple error page
func errorHandler(w http.ResponseWriter, statusCode int) {
	page := &Page{}
	w.WriteHeader(statusCode)
	txt := ""
	switch statusCode {
	case 400:
		txt = "400"
		page.Msg = "400 bad request"
	case 404:
		txt = "404"
		page.Msg = "404 page not found"
	case 405:
		txt = "405"
		page.Msg = "405 method not allowed"
	case 500:
		txt = "500"
		page.Msg = "500 internal server error"
	default:
		txt = "000"
		page.Msg = "an unexpected error occurred"
	}
	page.Title, _ = generator.GenArt(txt, "standard")
	errTmpl.Execute(w, page)
}

// checkRequired() checks a list of files if they exist or not.
// log.Fatal() if a file is not found
func checkRequired(reqFiles []string) {
	for _, file := range reqFiles {
		if !isFileThere(file) {
			log.Fatalf("Required file missing: %s\n", file)
		}
	}
}

// isFileThere() simply checks if a file exists
func isFileThere(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
