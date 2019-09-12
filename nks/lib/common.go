package lib

import (
	"encoding/json"
	"fmt"
)

var templates = map[string]SolutionTemplate{
	"jenkins":  jenkins,
}

func GetTemplateAsJson(s string) (string, error) {
	t := templates[s]
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("error while attempting to get solution template %s:\n\t:%v", s, err)
	}
	return string(b), err
}