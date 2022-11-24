package main

import (
	"context"

	"github.com/corezoid/gitcall-go-runner/gitcall"
	"github.com/rs/zerolog"
	"github.com/viplifes/simulator-cmd/command"
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
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	gitcall.Handle(usercode)
}
