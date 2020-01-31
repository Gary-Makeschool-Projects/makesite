package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

// structure for storing text file data
type text struct {
	Title   string
	Content string
}
// flag parser
func parser() (string, bool) {
	var path string
	var serve bool

	flag.StringVar(&path, "file", "", "The path to the text file we will convert.")
	flag.BoolVar(&serve, "", false, "Local hosting generated HTML file ")
	flag.Parse()

	if path != "" {
		return path, serve
	}

	flag.PrintDefaults()
	return path, serve
}

// function for reading file content
func read(path string) string {

	// parse path and read file content
	content, err := ioutil.ReadFile(path)
	// show error
	if err != nil {
		panic(err)
	}
	return string(content)

}
// function for creation of template from data
func createTemplate(content string) (*template.Template, data) {
	contentType := text{Content: content, Title:}

	tmpl, err := template.ParseFiles("index.tmpl")

	if err != nil {
		log.Fatal(err)
	}

	return tmpl, contentType
}
// function for rendering html 
func render(w http.ResponseWriter, r *http.Request) {
	data := text{"SSG", "Call of Duty"}

	fp := path.Join("templates", "index.html")

	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
// function for listening and serving file to localhost
func server() {
	http.HandleFunc("/", render)
	http.ListenAndServe(":5000", nil)
}
func main() {
	// parse all user flags
	path, serve := parser()
	// check if path is empty
	if path != "" {
		content := read(path)
		// create server
		if serve != false {
			server()
		}

	} else {
		// no specified path
		fmt.Print("You must provide a path to file with --file")
	}

}
