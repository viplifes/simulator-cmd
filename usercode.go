package main

import (
	"context"
	"github.com/corezoid/gitcall-examples/go/command"
)

func usercode(_ context.Context, data map[string]interface{}) error {

	cmd := data["cmd"].(string)

	params, err := command.Run(cmd, data)

	if err != nil {
		return err
	}

	data["cmdResult"] = params
	return nil
}

//func main() {
//	gitcall.Handle(usercode)
//}
