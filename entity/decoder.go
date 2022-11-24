package entity

import (
	"github.com/mitchellh/mapstructure"
)

func Decode(resp map[string]interface{}, response interface{}) {
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   response,
		TagName:  "json",
	})
	decoder.Decode(resp)
}
