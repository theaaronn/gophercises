package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	// Fields must be exported so Yaml package can access them
	Path string `yaml:"path"` // Corrected field name and YAML tag
	URL  string `yaml:"url"`  // Corrected field name and YAML tag
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
}

func buildMap(pathsUrls []pathUrl) map[string]string {
	pathMap := make(map[string]string, len(pathsUrls))
	for _, pu := range pathsUrls {
		pathMap[pu.Path] = pu.URL
	}
	return pathMap
}

func parseYAML(inbound []byte) ([]pathUrl, error) {
	var vessel []pathUrl
	err := yaml.Unmarshal(inbound, &vessel)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	if len(vessel) == 0 {
		return nil, fmt.Errorf("YAML is empty or contains no valid entries")
	}

	// Validate that the path and URL are not empty
	for _, entry := range vessel {
		if entry.Path == "" || entry.URL == "" {
			return nil, fmt.Errorf("YAML contains entries with empty path or url")
		}
	}

	return vessel, nil
}

func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathToUrls, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathsMap := buildMap(pathToUrls)
	return MapHandler(pathsMap, fallback), nil
}
