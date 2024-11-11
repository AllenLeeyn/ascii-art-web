package main

import (
	"ascii-art-web/pkg/generator"
	"net/http"
	"strings"
)

// homeHandler() handlers GET and POST request to "/" and "/ascii-art"
func homeHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" && req.URL.Path != "/ascii-art" {
		errorHandler(w, http.StatusNotFound)
		return
	}
	page := &AsciiPg{Banner: "standard"}

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

// handlePost() grabs user input using getFormInputs(),
// and clean up the text before passing it to GenArt().
// Uses the output of GenArt() to populate the response.
func handlePost(w http.ResponseWriter, req *http.Request, page *AsciiPg) {
	text, banner, formErr := getFormInputs(req)
	if formErr != "" {
		errorHandler(w, http.StatusBadRequest)
		return
	}
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.Trim(text, "\n")

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
	default:
		txt = "500"
		page.Msg = "500 internal server error"
	}
	page.Title, _ = generator.GenArt(txt, "standard")
	errTmpl.Execute(w, page)
}
