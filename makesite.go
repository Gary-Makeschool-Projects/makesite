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
	"path/filepath"

	"gopkg.in/russross/blackfriday.v2"
)

// structure for storing text file data
type text struct {
	Content string
}

// flag parser
func parser() (string, bool, string) {
	var path string
	var serve bool
	var dir string
	flag.StringVar(&path, "file", "", "The path to the text file we will convert.")
	flag.BoolVar(&serve, "", false, "Local hosting generated HTML file ")
	flag.StringVar(&dir, "dir", "", "Path to directory")
	flag.Parse()
	return path, serve, dir
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

	return buf.String()
}

func parseHTML(tmpl *template.Template, contentType text, fileName string) {
	template := writeHTML(tmpl, contentType)

	ioutil.WriteFile(fileName, []byte(template), 0666)
}

func convertMarkdown(path string) {
	content, err := ioutil.ReadFile(path)
	// nil check
	if err != nil {
		panic(err)
	}
	output := blackfriday.Run(content)
	name := getFileFromPath(path)
	ioutil.WriteFile(name+".html", output, 0666)
}

func getFileFromPath(path string) string {
	file := filepath.Base(path) // file name
	// get file extension
	extension := filepath.Ext(file)
	// remove the file extension
	name := file[0 : len(file)-len(extension)]
	// return the file name
	return name
}

func getFilesFromDirectory(path string) []string {
	allfiles := []string{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		extension := filepath.Ext(f.Name())
		if extension == ".txt" {
			allfiles = append(allfiles, path+"/"+f.Name())
		}

	}
	return allfiles
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
	path, serve, dir := parser()

	if path != "" {
		// check the file extension type
		extension := filepath.Ext(path)
		// get file name
		name := getFileFromPath(path)
		// if passed a txt file parse the txt file
		if extension == ".txt" {
			content := read(path)
			tmpl, data := createTemplate(content) // create the template
			parseHTML(tmpl, data, name+".html")
			// if passed a markdown file parse the markdown
		} else if extension == ".md" {
			convertMarkdown(path)
		} else {
			fmt.Print("cannot parse this file type")
		}

		// create server
		if serve != false {
			server()

		}

	} else if dir != "" {
		files := getFilesFromDirectory(dir) // list of all txt file in director
		ammount := len(files)               // amount of files in directory that will be generated
		for _, file := range files {
			name := getFileFromPath(file)
			content := read(file)
			tmpl, data := createTemplate(content) // create the template
			parseHTML(tmpl, data, name+".html")
		}
		fmt.Println("Success! ", ammount, " files were generated")

	} else {
		// no specified flag
		fmt.Print("You must provide a flag --file or -dir")
	}

}
