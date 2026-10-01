package input

import (
	"log"

	"github.com/goccy/go-yaml"
)

func ParseShopCategories(yamlData []byte) (map[string]string, error) {
	var categories map[string]string

	if err := yaml.Unmarshal(yamlData, &categories); err != nil {
		log.Printf("Error while unmarshalling YAML file: %v\n", err)
		return nil, err
	}

	return categories, nil
}
