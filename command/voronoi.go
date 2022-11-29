package command

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"

	"github.com/llgcode/draw2d/draw2dimg"

	"github.com/pzsz/voronoi"
	"github.com/viplifes/simulator-cmd/entity"
	"github.com/viplifes/simulator-cmd/simulator"
)

const voronoiNodeRef = "TestVoronoi"

func VoronoiAdd(data map[string]interface{}) (string, error) {

	token := data["token"].(string)
	layerId := data["layerId"].(string)
	formId := data["formId"].(string)
	simulatorAccId := data["simulatorAccId"].(string)

	client := simulator.New(token)

	resp, err := client.Get("graph_layers/"+layerId, nil)

	if err != nil {
		return "", err
	}

	var graph entity.LayerActors
	entity.Decode(resp, &graph)

	lineWidth := 2.0
	overflow := 100.0
	fixMinusX := 0.0
	fixMinusY := 0.0

	maxX := 0.0
	maxY := 0.0

	var nodes []entity.Actor

	for _, v := range graph.Data.Nodes {
		if v.Ref == voronoiNodeRef {
			//
		} else if v.Ref == "command-line" {
			//
		} else {
			nodes = append(nodes, v)
		}
	}

	for _, v := range nodes {

		// fix
		if v.Position.X < fixMinusX {
			fixMinusX = v.Position.X
		}
		if v.Position.Y < fixMinusY {
			fixMinusY = v.Position.Y
		}

		//x y
		if v.Position.X > maxX || maxX == 0.0 {
			maxX = v.Position.X
		}

		if v.Position.Y > maxY || maxY == 0.0 {
			maxY = v.Position.Y
		}
	}

	if fixMinusX < 0 {
		fixMinusX = fixMinusX - overflow
	}
	if fixMinusY < 0 {
		fixMinusY = fixMinusY - overflow
	}

	maxX = maxX + overflow - fixMinusX
	maxY = maxY + overflow - fixMinusY

	bbox := voronoi.NewBBox(0, maxX, 0, maxY)
	sites := make([]voronoi.Vertex, len(nodes))
	for j := 0; j < len(nodes); j++ {
		sites[j].X = nodes[j].Position.X - fixMinusX
		sites[j].Y = nodes[j].Position.Y - fixMinusY
	}

	diagram := voronoi.ComputeDiagram(sites, bbox, false)

	img := image.NewRGBA(image.Rect(0, 0, int(maxX)+int(lineWidth), int(maxY)+int(lineWidth)))

	l := draw2dimg.NewGraphicContext(img)
	l.SetLineWidth(lineWidth)

	// Iterate over cells
	for _, cell := range diagram.Cells {
		l.SetFillColor(color.Transparent)
		l.SetStrokeColor(color.RGBA{0xFF, 0x66, 0x00, 0xFF})
		for _, hedge := range cell.Halfedges {
			a := hedge.GetStartpoint()
			b := hedge.GetEndpoint()
			l.MoveTo(a.X, a.Y)
			l.LineTo(b.X, b.Y)
		}
		l.FillStroke()
	}

	l.Close()

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)

	err = png.Encode(foo, img)
	if err != nil {
		log.Println(err)
	}

	resp, _ = client.Upload("upload/"+simulatorAccId, b.Bytes(), nil)

	var file entity.UploadResponse
	entity.Decode(resp, &file)

	////// CREATE NEW
	resp, _ = client.Post("actors/actor/"+formId, simulator.Request{"ref": voronoiNodeRef, "appId": simulatorAccId, "data": simulator.Request{}, "title": "myTestVoronoi", "pictureObject": simulator.Request{
		"img":    file.Data.FileName,
		"width":  maxX,
		"height": maxY,
		"type":   "image",
	}}, nil)

	var actor entity.LayerActor
	entity.Decode(resp, &actor)

	////// ADD NEW TO LAYER
	client.Post("graph_layers/actors/"+layerId, []simulator.Request{
		{"action": "create", "data": simulator.Request{
			"id":   actor.Data.Id,
			"type": "node",
			"position": simulator.Request{
				"x": ((maxX + lineWidth*2) / 2) + fixMinusX,
				"y": ((maxY + lineWidth*2) / 2) + fixMinusY,
			}}},
	}, nil)

	//////

	return fmt.Sprintf("[ok] created Voronoi diagram for %d actors", len(graph.Data.Nodes)), nil
}

func VoronoiRemove(data map[string]interface{}) (string, error) {

	token := data["token"].(string)
	layerId := data["layerId"].(string)
	formId := data["formId"].(string)

	client := simulator.New(token)

	////// DELETE PREV ACTOR
	client.Delete("actors/ref/"+formId+"/"+voronoiNodeRef, nil, nil)

	resp, err := client.Get("graph_layers/"+layerId, nil)
	if err != nil {
		return "", err
	}

	var graph entity.LayerActors
	entity.Decode(resp, &graph)

	client.Put("graph_layers/actors/"+layerId, []simulator.Request{
		{
			"id": graph.Data.Nodes[0].LaId,
			"position": simulator.Request{
				"x": graph.Data.Nodes[0].Position.X + 0.001,
				"y": graph.Data.Nodes[0].Position.Y,
			}},
	}, nil)

	return fmt.Sprintf("[ok] removed Voronoi diagram"), nil
}
