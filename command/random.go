package command

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/AvraamMavridis/randomcolor"
	"github.com/corezoid/gitcall-examples/go/entity"
	"github.com/corezoid/gitcall-examples/go/simulator"
)

func RandomColor(data map[string]interface{}) (string, error) {

	token := data["token"].(string)
	layerId := data["layerId"].(string)

	client := simulator.New(token)

	resp, err := client.Get("graph_layers/"+layerId, nil)

	if err != nil {
		return "", err
	}

	var graph entity.LayerActors
	entity.Decode(resp, &graph)

	var wg sync.WaitGroup
	for _, v := range graph.Data.Nodes {
		color := randomcolor.GetRandomColorInHex()
		wg.Add(1)
		go func(v entity.Actor, client *simulator.Client, token string, color string) {
			defer wg.Done()
			client.Put("actors/actor/"+strconv.Itoa(v.FormId)+"/"+v.Id, simulator.Request{"color": color}, simulator.Request{"replaceEmpty": "false"})
		}(v, client, token, color)

	}
	wg.Wait()
	return fmt.Sprintf("[ok] updated color for %d actors", len(graph.Data.Nodes)), nil
}
