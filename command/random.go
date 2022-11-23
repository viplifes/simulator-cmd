package command

import (
	"github.com/AvraamMavridis/randomcolor"
	"github.com/corezoid/gitcall-examples/go/entity"
	"github.com/corezoid/gitcall-examples/go/simulator"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"sync"
)

func RandomColor(data map[string]interface{}) (map[string]interface{}, error) {

	token := data["token"].(string)
	layerId := data["layerId"].(string)

	client := simulator.New(token)

	resp, err := client.Get("graph_layers/"+layerId, nil)

	var response entity.LayerActors
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &response,
		TagName:  "json",
	})
	decoder.Decode(resp)
	var wg sync.WaitGroup
	for _, v := range response.Data.Nodes {
		color := randomcolor.GetRandomColorInHex()
		wg.Add(1)
		go func(v entity.Actor, client *simulator.Client, token string, color string) {
			defer wg.Done()
			client.Put("actors/actor/"+strconv.Itoa(v.FormId)+"/"+v.Id, simulator.Request{"color": color}, simulator.Request{"replaceEmpty": "false"})
		}(v, client, token, color)

	}
	wg.Wait()
	return resp, err
}
