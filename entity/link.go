package entity

type LinkActor struct {
	Data LinkActorItem `json:"data"`
}

type LinkActorItem struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Source        string `json:"source"`
	Target        string `json:"target"`
	EdgeType      string `json:"edgeType"`
	LinkedActorId string `json:"linkedActorId"`
}
