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

var aboutTmpl = template.Must(template.ParseFiles("templates/about.html"))
var errTmpl = template.Must(template.ParseFiles("templates/error.html"))
var indexTmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	generator.GetStyles()
	http.Handle("/static/", http.FileServer(http.Dir("assets/")))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
