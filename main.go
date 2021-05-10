package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/lucas13004/groupie-tackers/handler"
)

func main() {

	templates, err := template.ParseGlob("templates/*")
	if err != nil {
		panic(err)
	}

	sh := handler.New(templates)
	http.HandleFunc("/", sh.Route)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
