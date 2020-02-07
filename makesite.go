package main

import (
	"bytes"
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
	Content string
}

// flag parser
func parser() (string, bool) {
	var path string
	var serve bool
	flag.StringVar(&path, "file", "", "The path to the text file we will convert.")
	flag.BoolVar(&serve, "", false, "Local hosting generated HTML file ")
	flag.Parse()
	// check if a file path was given in the flag
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
	// nil check
	if err != nil {
		panic(err)
	}
	// return the contents of the file
	return string(content)

}

// function for creation of template from data
func createTemplate(content string) (*template.Template, text) {
	contentType := text{Content: content}

	tmpl, err := template.ParseFiles("templates/index.tmpl")
	// nil check
	if err != nil {
		log.Fatal(err)
	}
	// return the template and the content type
	return tmpl, contentType
}

func writeHTML(tmpl *template.Template, contentType text) string {
	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, contentType); err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())
	return buf.String()
}

func parseHTML(tmpl *template.Template, contentType text, fileName string) {
	template := writeHTML(tmpl, contentType)

	ioutil.WriteFile(fileName, []byte(template), 0666)
}

// function for rendering html
func render(w http.ResponseWriter, r *http.Request) {
	data := text{"SSG"}

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
		tmpl, data := createTemplate(content) // create the template
		fmt.Print(data)
		fmt.Print(tmpl)
		parseHTML(tmpl, data, "generated.html")
		// create server
		if serve != false {
			server()
		}

	} else {
		// no specified path
		fmt.Print("You must provide a path to file with --file")
	}

}
