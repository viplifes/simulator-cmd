package command

import (
	"fmt"
	"sync"

	"github.com/viplifes/simulator-cmd/command/tsp"
	"github.com/viplifes/simulator-cmd/entity"
	"github.com/viplifes/simulator-cmd/simulator"
)

const EdgeTestName = "CmdTestEdge"

func EdgesAdd(data map[string]interface{}) (string, error) {
	token := data["token"].(string)
	layerId := data["layerId"].(string)
	simulatorAccId := data["simulatorAccId"].(string)
	client := simulator.New(token)

	resp, _ := client.Get("graph_layers/"+layerId, nil)

	var graph entity.LayerActors
	entity.Decode(resp, &graph)

	var nodes []entity.Actor
	for _, v := range graph.Data.Nodes {
		if v.Ref != "command-line" {
			nodes = append(nodes, v)
		}
	}

	/// TSP
	nodesNew := tsp.TspRun(nodes, len(nodes)*150)
	//	nodesNew := TspRun(nodes, 1000)
	var links []entity.LinkActor

	for i := 0; i < len(nodesNew)-1; i++ {

		source := nodesNew[i]
		target := nodesNew[i+1]

		resp, _ = client.Post("actors/link/"+simulatorAccId, simulator.Request{"source": source.Id, "target": target.Id, "edgeTypeId": 13, "name": EdgeTestName}, nil)

		var link entity.LinkActor
		entity.Decode(resp, &link)

		found := false
		for _, v := range graph.Data.Edges {
			if link.Data.Id == v.Id {
				found = true
			}
		}

		if !found {
			links = append(links, link)
		}

	}

	var reqManage []simulator.Request

	for _, v := range links {
		reqManage = append(reqManage, simulator.Request{"action": "create", "data": simulator.Request{"id": v.Data.Id, "type": "edge"}})
	}

	if len(reqManage) > 0 {
		client.Post("graph_layers/actors/"+layerId, reqManage, nil)
	}

	return fmt.Sprintf("[ok] created %d edges", len(reqManage)), nil
}

func EdgesRemove(data map[string]interface{}) (string, error) {

	token := data["token"].(string)
	layerId := data["layerId"].(string)
	client := simulator.New(token)

	resp, _ := client.Get("graph_layers/"+layerId, nil)

	var graph entity.LayerActors
	entity.Decode(resp, &graph)

	count := 0
	var wg sync.WaitGroup
	for _, v := range graph.Data.Edges {
		if v.Name == EdgeTestName {
			count++
			wg.Add(1)
			go func(v entity.LinkActorItem, client *simulator.Client, token string) {
				defer wg.Done()
				client.Delete("actors/link/"+v.Id, nil, nil)
			}(v, client, token)
		}
	}
	wg.Wait()

	return fmt.Sprintf("[ok] removed %d edges", count), nil
}
