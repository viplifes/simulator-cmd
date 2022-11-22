package command

import (
	"errors"
	"strings"
)

func Run(cmd string, data map[string]interface{}) (map[string]interface{}, error) {
	words := strings.Fields(cmd)
	if words[0] == "random" && words[1] == "color" {
		return RandomColor(data)
	} else if words[0] == "tsp" {
		return Tsp(data)
	} else {
		return nil, errors.New("command '" + cmd + "' not_found")
	}
}
