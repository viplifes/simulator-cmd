package main

import (
	"context"

	"github.com/corezoid/gitcall-examples/go/command"
	"github.com/corezoid/gitcall-go-runner/gitcall"
)

func usercode(_ context.Context, data map[string]interface{}) error {
	cmd := data["cmd"].(string)
	result, err := command.Run(cmd, data)
	if err != nil {
		return err
	}
	data["cmdResult"] = result
	return nil
}

func main() {
	gitcall.Handle(usercode)
}
