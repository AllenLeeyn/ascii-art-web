package main

import (
	"ascii-art-web/pkg/generator"
	"html/template"
	"log"
	"net/http"
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
	generator.GetStyles()
	http.Handle("GET /static/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", homeHandler)

	log.Println("Listening and serving on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
