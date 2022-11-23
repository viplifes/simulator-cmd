package entity

type LayerActors struct {
	Data LayerActorsList `json:"data"`
}

type LayerActorsList struct {
	Nodes []Actor `json:"nodes"`
}

type Actor struct {
	Id       string        `json:"id"`
	Title    string        `json:"title"`
	Position ActorPosition `json:"position"`
	FormId   int           `json:"formId"`
}

type ActorPosition struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}
