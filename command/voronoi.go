package command

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/viplifes/simulator-cmd/command/voronoi"
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

	maxX := 1000.0
	maxY := 1000.0

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
		if v.Position.X > maxX {
			maxX = v.Position.X
		}
		if v.Position.Y > maxY {
			maxY = v.Position.Y
		}
	}

	vor := voronoi.Voronoi{}
	svg := SVG{
		Title:       "Voronoi diagram",
		Description: "Edges and points",
		Width:       maxX,
		Height:      maxY,
		StrokeWidth: 3,
		PointRadius: 5,
		Vertices:    make([]*voronoi.Point, len(nodes)),
	}
	for i, v := range nodes {
		svg.Vertices[i] = voronoi.Pt(v.Position.X, v.Position.Y)
	}
	svg.Edges = vor.GetEdges(&svg.Vertices, svg.Width, svg.Height)
	tmpl := template.Must(template.New("svg").Parse(TEMPLATE))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, svg); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	////// DELETE PREV ACTOR
	client.Delete("actors/ref/"+formId+"/"+voronoiNodeRef, nil, nil)

	////// UPLOAD FILE
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	sVGtoPNG(buf.Bytes(), foo)

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
				"x": maxX / 2,
				"y": maxY / 2,
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

const TEMPLATE = `<?xml version="1.0" ?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN"
  "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="{{.Width}}px" height="{{.Height}}px" viewBox="0 0 {{.Width}} {{.Height}}"
     xmlns="http://www.w3.org/2000/svg" version="1.1">
  <title>{{.Title}}</title>
  <desc>{{.Description}}</desc>
  <!-- Edges -->
  <g stroke="red" stroke-width="{{.StrokeWidth}}" fill="none">
    {{range .Edges}}<path d="M{{.Start.X}},{{.Start.Y}} L{{.End.X}},{{.End.Y}}" />
    {{end}}</g>
</svg>`

type SVG struct {
	Width       float64
	Height      float64
	Edges       voronoi.Edges
	Vertices    voronoi.Vertices
	Title       string
	Description string
	StrokeWidth float64
	PointRadius float64
}

func sVGtoPNG(svg []byte, pngWriter io.Writer) {
	icon, _ := oksvg.ReadIconStream(bytes.NewBuffer(svg))
	w := int(icon.ViewBox.W)
	h := int(icon.ViewBox.H)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	err := png.Encode(pngWriter, rgba)
	if err != nil {
		log.Println(err)
	}
}
