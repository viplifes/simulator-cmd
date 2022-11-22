package command

import (
	"github.com/AvraamMavridis/randomcolor"
	"github.com/corezoid/gitcall-examples/go/entity"
	"github.com/corezoid/gitcall-examples/go/simulator"
	"github.com/mitchellh/mapstructure"
	"sync"
)

func RandomColor(data map[string]interface{}) (map[string]interface{}, error) {

	token := data["token"].(string)
	layerId := data["layerId"].(string)
	formId := data["formId"].(string)

	client := simulator.New(token)

	resp, err := client.Get("layer_actors_filters/"+layerId+"/"+formId, simulator.Request{"filter": "id,title,position"})

	var response entity.LayerActors
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &response,
		TagName:  "json",
	})
	decoder.Decode(resp)
	var wg sync.WaitGroup
	for _, v := range response.Data.List {
		color := randomcolor.GetRandomColorInHex()
		wg.Add(1)
		go func(v entity.Actor, client *simulator.Client, token string, color string, formId string) {
			client.Put("actors/actor/"+formId+"/"+v.Id, simulator.Request{"color": color}, simulator.Request{"replaceEmpty": "false"})
			defer wg.Done()
		}(v, client, token, color, formId)

	}
	wg.Wait()
	return resp, err
}
