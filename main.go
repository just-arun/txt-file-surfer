package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "bar", "a string var")
	flag.Parse()

	fmt.Println(path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("started")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dataLi := `
		<ul>
		`
		for _, f := range files {
			fmt.Println(f.Name(), f.IsDir())
			dataLi += `<li>
			<a href="/file?filePath=` + path + "/" + f.Name() + `">` + f.Name() + `</a>
			</li>`
		}
		dataLi += `</ul>`
		t, err := template.ParseFiles("./index.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := t.Execute(w, template.HTML(dataLi)); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Query().Get("filePath")
		file, _ := os.Open(filePath)
		fmt.Println(filePath)
		scanner := bufio.NewScanner(file)
		var fileContent string
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			fileContent = scanner.Text()
		}
		t, err := template.ParseFiles("./file.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := t.Execute(w, fileContent); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Error(w, path, 500)
	})
	http.ListenAndServe(":9000", nil)
}
