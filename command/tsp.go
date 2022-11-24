package command

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/viplifes/simulator-cmd/command/tsp/base"
	ga "github.com/viplifes/simulator-cmd/command/tsp/geneticAlgorithm"
	"github.com/viplifes/simulator-cmd/entity"
)

func TspRun(nodes []entity.Actor, gen int) []entity.Actor {
	rand.Seed(time.Now().UTC().UnixNano())

	// Init TourManager
	tm := base.TourManager{}
	tm.NewTourManager()
	// Add cities to TourManager
	for _, v := range nodes {
		tm.AddCity(base.GenerateCity(int(v.Position.X), int(v.Position.Y), v))
	}
	result := tspGA(&tm, gen)

	var nodesNew []entity.Actor
	for i := 0; i < len(result); i++ {
		nodesNew = append(nodesNew, result[i])
	}

	return nodesNew

}

// tspGA : Travelling sales person with genetic algorithm
// input :- TourManager, Number of generations
func tspGA(tm *base.TourManager, gen int) []entity.Actor {
	p := base.Population{}
	// Population Size
	p.InitPopulation(3, *tm)

	// Get initial fittest tour and it's tour distance
	fmt.Println("Start....")
	iFit := p.GetFittest()

	// Map to store fittest tours
	fittestTours := make([]base.Tour, 0, gen+1)
	fittestTours = append(fittestTours, *iFit)
	// Evolve population "gen" number of times
	for i := 1; i < gen+1; i++ {
		p = ga.EvolvePopulation(p)
		fittestTours = append(fittestTours, *p.GetFittest())
	}
	return tourToPoints(p.GetFittest())
}

func tourToPoints(t *base.Tour) []entity.Actor {
	tLen := t.TourSize()
	nodes := make([]entity.Actor, tLen+1)

	c0 := t.GetCity(0)
	nodes[0] = c0.Actor
	nodes[tLen] = c0.Actor
	for i := 1; i < tLen; i++ {
		c := t.GetCity(i)
		nodes[i] = c.Actor
	}
	return nodes
}
