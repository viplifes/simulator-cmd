package base

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/viplifes/simulator-cmd/entity"
)

// City : coordinates of city
type City struct {
	x     int
	y     int
	Actor entity.Actor
}

// GenerateCity : Generate city with user defined coordinates
func GenerateCity(actor entity.Actor) City {
	c := City{}
	c.x = int(actor.Position.X)
	c.y = int(actor.Position.Y)
	c.Actor = actor
	return c
}

// SetLocation : User defined coordinates for a city
func (a *City) SetLocation(x int, y int) {
	a.x = x
	a.y = y
}

// DistanceTo : distance of current city to target city
func (a *City) DistanceTo(c City) float64 {
	idx := a.x - c.x
	idy := a.y - c.y

	if idx < 0 {
		idx = -idx
	}
	if idy < 0 {
		idy = -idy
	}

	fdx := float64(idx)
	fdy := float64(idy)

	fd := math.Sqrt((fdx * fdx) + (fdy * fdy))
	return fd
}

func (a *City) X() int {
	return a.x
}

func (a *City) Y() int {
	return a.y
}

func (a City) String() string {
	return fmt.Sprintf("{x%d y%d}", a.x, a.y)
}

// ShuffleCities : return a shuffled []City given input []City
func ShuffleCities(in []City) []City {
	out := make([]City, len(in), cap(in))
	perm := rand.Perm(len(in))
	for i, v := range perm {
		out[v] = in[i]
	}
	return out
}
