package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

type favorite struct {
	Color   string
	Game    string
	Content string
}

func parser(path string) string {
	// parse path
	content := []string{}
	file, err := os.Open(path)
	// check to see if file is empty
	if err != nil {
		panic(err)
	}
	// close file when operation done
	defer file.Close()
	//
	buf := make([]byte, 1000)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}

		content = append(content, string(buf[:n]))
	}

	return strings.Join(content, " ")
}
func createFile(content string) []byte {
	os.Create("/Users/ghost/go/src/github.com/imthaghost/ssg/example.html")
	b := []byte(content)

	return b

}

func render(w http.ResponseWriter, r *http.Request) {

	data := favorite{"Blue", "Call of Duty", "empty"}

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
func serve() {
	http.HandleFunc("/", render)
	http.ListenAndServe(":5000", nil)
}
func main() {

	var content string

	flag.StringVar(&content, "content", "", "Usage")

	flag.Parse()

	fmt.Println(content)
}
