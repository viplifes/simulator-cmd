package entity

type LayerActors struct {
	Data LayerActorsList `json:"data"`
}

type LayerActorsList struct {
	List  []Actor `json:"list"`
	Total int     `json:"total"`
}

type Actor struct {
	Id       string        `json:"id"`
	Title    string        `json:"title"`
	Position ActorPosition `json:"position"`
}

type ActorPosition struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}
