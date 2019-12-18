package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	str := `<html>
<head><title>hello, go</title></head>
<body><h1>Hello, Go~</h1></body>
</html>`
	w.Write([]byte(str))
}

func headers(w http.ResponseWriter, r *http.Request) {
	h := r.Header.Get("Accept-Encoding")
	fmt.Fprintln(w, h)
}

func body(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	fmt.Fprint(w, string(body))
}

func process(w http.ResponseWriter, r *http.Request) {
	// handle upload file
	r.ParseMultipartForm(1024)
	fileHeader := r.MultipartForm.File["upload"][0]
	if file, err := fileHeader.Open(); err == nil {
		if data, err := ioutil.ReadAll(file); err == nil {
			fmt.Fprintf(w, string(data))
		}
	}
}

type Post struct {
	User    string
	Threads []string
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	post := &Post{"heiheihei", []string{"one", "two", "three"}}
	json, _ := json.Marshal(post)
	w.Write(json)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/", index)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/body", body)
	http.HandleFunc("/process", process)
	http.HandleFunc("/json", jsonHandler)
	server.ListenAndServe()
}
