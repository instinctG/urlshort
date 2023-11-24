package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/instinctG/urlshort/handler
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := handler.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	//yamlData, err := getFile("data/map.yaml")
	jsonData, err := getFile("data/map.json")

	//yamlHandler, err := handler.YAMLHandler(yamlData, mapHandler)
	jsonHandler, err := handler.JSONHandler(jsonData, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func getFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dataFromFile, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return dataFromFile, err
}
