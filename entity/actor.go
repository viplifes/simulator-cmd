package entity

type LayerActors struct {
	Data LayerActorsList `json:"data"`
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
}

type ActorPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
