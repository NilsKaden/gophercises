package urlshort

import (
	"net/http"
	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// define the returned handler func
	return func(w http.ResponseWriter, r *http.Request) {
		// get path from Request
		path := r.URL.Path

		// check if key (path from request) exits in map
		if dest, ok := pathsToUrls[path]; ok {
			// if so, create a redirect to the value of the map (dest)
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		// serve fallback
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// fill variables with yaml data
	parsedYaml, err := parseYAML(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	// unmarshal references the struct for mapping yaml to variables
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}

	return pathUrls, nil
}


func buildMap(pathUrls []pathUrl) map[string]string {
	// make preallocates the space required for the map. Additionally, it supports maps with len != cap
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}

	return pathsToUrls
}

// interface for mapping yaml data to variables 
type pathUrl struct {
	Path string `yaml:"path"`
	URL string `yaml:"url"`
}

