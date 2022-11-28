package entity

type LayerActors struct {
	Data LayerActorsList `json:"data"`
}

type LayerActor struct {
	Data Actor `json:"data"`
}

type UploadResponse struct {
	Data struct {
		FileName string `json:"fileName"`
	}
}

type LayerActorsList struct {
	Nodes []Actor         `json:"nodes"`
	Edges []LinkActorItem `json:"edges"`
}

type Actor struct {
	Ref      string        `json:"ref"`
	Id       string        `json:"id"`
	Title    string        `json:"title"`
	Position ActorPosition `json:"position"`
	FormId   int           `json:"formId"`
	LaId     int           `json:"laId"`
}

type ActorPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
