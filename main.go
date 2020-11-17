package main

import (
	"fmt"
	"github.com/bradfitz/slice"
	"io/ioutil"
	"log"
	"math"
	"strings"
)

type Tile struct {
	X            int
	Y            int
	distance     int
	cost         int
	CostDistance int
	Parent       interface{}
}

var maxX = 12
var maxY = 10

func (t *Tile) setCostDistance() {
	t.CostDistance = t.cost + t.distance
}

func (t *Tile) setDistance(targetX, targetY int) {
	t.distance = int(math.Abs(float64(targetX-t.X)) + math.Abs(float64(targetY-t.Y)))
}

func getPossibleTiles(gameMap []string, currentTile, targetTile *Tile) []*Tile {
	tilesDir := []*Tile{
		&Tile{X: currentTile.X, Y: currentTile.Y - 1, cost: currentTile.cost + 1, Parent: currentTile},
		&Tile{X: currentTile.X, Y: currentTile.Y + 1, cost: currentTile.cost + 1, Parent: currentTile},
		&Tile{X: currentTile.X - 1, Y: currentTile.Y, cost: currentTile.cost + 1, Parent: currentTile},
		&Tile{X: currentTile.X + 1, Y: currentTile.Y, cost: currentTile.cost + 1, Parent: currentTile},
	}

	possible := []*Tile{}

	for _, tile := range tilesDir {
		tile.setDistance(targetTile.X, targetTile.Y)
		if (tile.X >= 0 && tile.X <= maxX) && (tile.Y >= 0 && tile.Y <= maxY) {
			if gameMap[tile.Y][tile.X] == '.' || gameMap[tile.Y][tile.X] == ';' {
				possible = append(possible, tile)
			}
		}
	}

	return possible
}

func remove(s []*Tile, target *Tile) (res []*Tile) {
	for _, k := range s {
		if k != target {
			res = append(res, k)
		}
	}
	return
}

func main() {
	defer fmt.Println()
	fileMap, err := ioutil.ReadFile("map1.txt")
	if err != nil {
		log.Println(err)
	}
	mapStr := ""
	for _, k := range fileMap {
		mapStr += string(k)
	}
	// fmt.Print(mapStr)
	mapArr := strings.Split(mapStr, "\n")

	start := &Tile{
		X: 0,
		Y: 0,
	}

	finish := &Tile{
		X: 12,
		Y: 10,
	}

	start.setDistance(finish.X, finish.Y)

	activeTiles := []*Tile{start}
	visitedTiles := []*Tile{}

	for len(activeTiles) > 0 {

		for _, k := range activeTiles {
			fmt.Println(k)
		}

		slice.Sort(activeTiles, func(i, j int) bool {
			return activeTiles[i].CostDistance < activeTiles[j].CostDistance
		})

		var checkTile = activeTiles[0]

		if checkTile.X == finish.X && checkTile.Y == finish.Y {
			fmt.Println("EBAT NAHUI, POLUCHILOS' BLYAT'")
			return
		}

		visitedTiles = append(visitedTiles, checkTile)
		activeTiles = remove(activeTiles, checkTile)

		possible := getPossibleTiles(mapArr, checkTile, finish)
	Loop:
		for _, walkableTile := range possible {

			for _, k := range visitedTiles {
				if k.X == walkableTile.X && k.Y == walkableTile.Y {
					continue Loop
				}
			}

			for _, k := range activeTiles {
				if k.X == walkableTile.X && k.Y == walkableTile.Y {
					existingTile := k
					if existingTile.CostDistance > checkTile.CostDistance {
						activeTiles = remove(activeTiles, existingTile)
						continue Loop
					}
				}
			}
			activeTiles = append(activeTiles, walkableTile)
		}

	}
}
