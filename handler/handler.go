package handler

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"net/http"
)

type PathToUrl struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

// MapHandler will return a http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if val, ok := pathsToUrls[r.URL.String()]; ok {
			fmt.Println("Matched:", val)
			http.Redirect(w, r, val, http.StatusFound)
			return
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// a http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yamlData)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	fmt.Println(pathMap)
	return MapHandler(pathMap, fallback), nil
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(jsonData)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

func parseYaml(yamlData []byte) (pairs []PathToUrl, err error) {
	err = yaml.Unmarshal(yamlData, &pairs)
	return pairs, err
}

func buildMap(pairs []PathToUrl) map[string]string {
	builtMap := make(map[string]string)

	for _, p := range pairs {
		builtMap[p.Path] = p.Url

	}

	return builtMap
}

func parseJSON(jsonData []byte) (jsonPairs []PathToUrl, err error) {
	err = json.Unmarshal(jsonData, &jsonPairs)
	return
}
