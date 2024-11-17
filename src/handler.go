package main

import (
	"ascii-art-web/pkg/generator"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// aboutHandler handles GET requests to "/about"
func aboutHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		errorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	page := &Page{Title: "About"}
	aboutTmpl.Execute(w, page)
}

// homeHandler() handlers GET and POST request to "/" and "/ascii-art"
func homeHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" && req.URL.Path != "/ascii-art" {
		errorHandler(w, http.StatusNotFound)
		return
	}
	page := &AsciiPg{Banner: "standard"}

	page.P.Title, _ = generator.GenArt("ASCII-ART", "standard")
	page.P.Msg, _ = generator.GenArt("Enter text\nabove", "standard")

	switch req.Method {
	case "GET":
		indexTmpl.Execute(w, page)
	case "POST":
		handlePost(w, req, page)
	default:
		errorHandler(w, http.StatusMethodNotAllowed)
	}
}

// handlePost() grabs user input using getFormInputs().
// Uses the output of GenArt() to populate the response.
func handlePost(w http.ResponseWriter, req *http.Request, page *AsciiPg) {
	text, banner, formErr := getFormInputs(req)

	if formErr != "" {
		errorHandler(w, http.StatusBadRequest)
		return
	}
	if _, err := os.Stat("./assets/styles/" + banner + ".txt"); err != nil {
		errorHandler(w, http.StatusInternalServerError)
		return
	}
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.Trim(text, "\n")
	if art, err := generator.GenArt(text, banner); err != nil {
		page.P.Msg = "Failed to generate ASCII art: " + err.Error()
	} else {
		page.P.Msg = art
	}
	page.Text, page.Banner = text, banner

	indexTmpl.Execute(w, page)
}

// getFormInputs() gets text and banner input from the form in the POST request.
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
func errorHandler(w http.ResponseWriter, statusCode int) {
	page := &Page{Title: strconv.Itoa(statusCode)}
	w.WriteHeader(statusCode)
	switch page.Title {
	case "400":
		page.Msg = "400 bad request"
	case "404":
		page.Msg = "404 page not found"
	case "405":
		page.Msg = "405 method not allowed"
	default:
		page.Msg = "500 internal server error"
	}
	page.Title, _ = generator.GenArt(page.Title, "standard")
	errTmpl.Execute(w, page)
}
