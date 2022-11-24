package command

import (
	"errors"
	"strings"
)

func Run(cmd string, data map[string]interface{}) (string, error) {
	words := strings.Fields(cmd)
	if words[0] == "random" && words[1] == "color" {
		return RandomColor(data)
	} else if words[0] == "edges" && words[1] == "add" {
		return EdgesAdd(data)
	} else if words[0] == "edges" && words[1] == "remove" {
		return EdgesRemove(data)
	} else {
		return "", errors.New("[error] command '" + cmd + "' not found")
	}
}
