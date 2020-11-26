package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sort"
	"strings"
)

type Coords struct {
	x int
	y int
}

type Entity struct {
	name string
	id   int
	Coords
	param1 int
	param2 int
}

type Tile struct {
	Coords
	direction    string
	distance     int
	cost         int
	CostDistance int
	path  		 []string
	retreatPath  []string
	gameMap 	 []string
	isBox        bool
	Parent       *Tile
}

const (
	maxX       = 12
	maxY       = 10
	bombRadius = 2
)

func (t *Tile) setCostDistance() {
	t.CostDistance = t.cost + t.distance
}

func (t *Tile) setDistance(targetX, targetY int) {
	t.distance = int(math.Abs(float64(targetX-t.x)) + math.Abs(float64(targetY-t.y)))
}

func getResult(checkTile *Tile) (path []string) {
	fmt.Println("\nMap:")
	for _, k := range checkTile.gameMap {
		fmt.Println(k)
	}
	fmt.Println()

	return
}

func remove(s []*Tile, target *Tile) (res []*Tile) {
	for _, k := range s {
		if k != target {
			res = append(res, k)
		}
	}
	return
}

func sortTiles(activeTiles []*Tile) {
	sort.Slice(activeTiles, func(i, j int) bool {
		return activeTiles[i].CostDistance < activeTiles[j].CostDistance
	})
}

func makeMapArray() []string {
	fileMap, err := ioutil.ReadFile("map1.txt")
	if err != nil {
		log.Println(err)
	}
	mapStr := ""
	for _, k := range fileMap {
		mapStr += string(k)
	}

	return strings.Split(mapStr, "\n")
}

func reverse(arr []string) (res []string) {
	opposites := map[string]string{
		"left":  "right",
		"right": "left",
		"up":    "down",
		"down":  "up",
	}
	for i := range arr {
		if opposites[arr[i]] != "" {
			res = append([]string{opposites[arr[i]]}, res...)
		}
	}
	return res
}

func getPossibleTiles(currentTile, targetTile *Tile, visitedTiles []*Tile) []*Tile {

	mapInst := make([]string, maxY+1)

	for i, k := range currentTile.gameMap {
		mapInst[i] = k
	}

	// initialize 4 direction
	tilesDir := []*Tile{
		&Tile{Coords: Coords{x: currentTile.x, y: currentTile.y - 1}, cost: currentTile.cost + 1, Parent: currentTile, direction: "up", gameMap: mapInst},
		&Tile{Coords: Coords{x: currentTile.x, y: currentTile.y + 1}, cost: currentTile.cost + 1, Parent: currentTile, direction: "down", gameMap: mapInst},
		&Tile{Coords: Coords{x: currentTile.x - 1, y: currentTile.y}, cost: currentTile.cost + 1, Parent: currentTile, direction: "left", gameMap: mapInst},
		&Tile{Coords: Coords{x: currentTile.x + 1, y: currentTile.y}, cost: currentTile.cost + 1, Parent: currentTile, direction: "right", gameMap: mapInst},
	}

	possible := []*Tile{}

	// initialize possible directions
	for _, tile := range tilesDir {
		tile.gameMap = make([]string, maxY+1)

		for i, k := range currentTile.gameMap {
			tile.gameMap[i] = k
		}

		tile.setDistance(targetTile.x, targetTile.y)
		if (tile.x >= 0 && tile.x <= maxX) && (tile.y >= 0 && tile.y <= maxY) {
			if tile.gameMap[tile.y][tile.x] != '!' {

				// meeting breakable box (spends 11-12 ticks, depending on position)
				if tile.gameMap[tile.y][tile.x] == ';' {
					tile.updateMap('.')
					for _, k := range tile.gameMap {
						fmt.Println(k)
					}
					
					fmt.Println()
					
					tile.setRetreatPath()
					tile.isBox = true
					tile.cost += 10
				}

				tile.setCostDistance()
				possible = append(possible, tile)
			}
		}
	}

	return possible
}

func (t *Tile) setRetreatPath() {

	retreatCoord := bfs(t.gameMap, Coords{x: t.Parent.x, y: t.Parent.y})
	fmt.Println("Starting Coordinates: ", t.Parent.Coords)
	fmt.Println("Retreat Coordinates: ", retreatCoord)
	return
	retreatPath := AStar(t.gameMap, t.Parent.Coords, retreatCoord)
	fmt.Println("Retreat Path: ", retreatPath)
	return 
	retreat := append([]string{"bomb"}, retreatPath...)
	retreat = append(retreat, []string{"stay", "stay"}...)
	retreat = append(retreat, reverse(retreat)...)
	retreat = append(retreat, t.direction)
	t.retreatPath = append(t.retreatPath, retreat...)

}

func replaceAtIndex(in string, r rune, i int) string {
    out := []rune(in)
    out[i] = r
    return string(out)
}

func (t *Tile) updateMap(r rune) {
	t.gameMap[t.y] = replaceAtIndex(t.gameMap[t.y], r, t.x)
}

func AStar(mapArr []string, startCoords, finishCoords Coords) []string {
	checkTile := &Tile{}

	start := &Tile{
		Coords: startCoords,
		gameMap: mapArr,
	}

	finish := &Tile{
		Coords: finishCoords,
	}

	start.setDistance(finish.x, finish.y)

	fmt.Println(finish)
	return []string{}

	activeTiles := []*Tile{start}
	visitedTiles := []*Tile{}

	for len(activeTiles) > 0 {

		// Sorting tiles in stack to chose better option
		sortTiles(activeTiles)

		// Best option
		checkTile = activeTiles[0]

		if checkTile.Parent != nil {
			if !checkTile.isBox {
				checkTile.path = append(checkTile.Parent.path, checkTile.direction)
			}
		}		
		
		if checkTile.x == finish.x && checkTile.y == finish.y {
			getResult(checkTile)
			return checkTile.path
		}

		visitedTiles = append(visitedTiles, checkTile)
		activeTiles = remove(activeTiles, checkTile)

		possible := getPossibleTiles(checkTile, finish, visitedTiles)

Loop:
		for _, walkableTile := range possible {

			// Prevent entering visited cell
			for _, k := range visitedTiles {
				if k.x == walkableTile.x && k.y == walkableTile.y {
					continue Loop
				}
			}

			activeTiles = append(activeTiles, walkableTile)
		}
	}

	// path = getResult(checkTile, mapArr)

	fmt.Println(checkTile.path)
	return checkTile.path
}

func main() {
	mapArr := makeMapArray()
	AStar(mapArr, Coords{x: 0, y: 0}, Coords{x: maxX, y: maxY})
}

func explosionArea(gameMap []string, bombC Coords) map[Coords]bool {
	explosionArea := map[Coords]bool{
		bombC: true,
	}
	for bCounter := 1; bCounter <= bombRadius; bCounter++ {
		explCell := Coords{bombC.x - bCounter, bombC.y}
		if 0 <= explCell.x && gameMap[explCell.y][explCell.x] == '.' {
			explosionArea[explCell] = true
		} else {
			break
		}
	}

	for bCounter := 1; bCounter <= bombRadius; bCounter++ {
		explCell := Coords{bombC.x + bCounter, bombC.y}
		if explCell.x <= maxX && gameMap[explCell.y][explCell.x] == '.' {
			explosionArea[explCell] = true
		} else {
			break
		}
	}

	for bCounter := 1; bCounter <= bombRadius; bCounter++ {
		explCell := Coords{bombC.x, bombC.y - bCounter}
		if 0 <= explCell.y && gameMap[explCell.y][explCell.x] == '.' {
			explosionArea[explCell] = true
		} else {
			break
		}
	}

	for bCounter := 1; bCounter <= bombRadius; bCounter++ {
		explCell := Coords{bombC.x, bombC.y + bCounter}
		if explCell.y <= maxY && gameMap[explCell.y][explCell.x] == '.' {
			explosionArea[explCell] = true
		} else {
			break
		}
	}

	return explosionArea
}

func bfs(gameMap []string, start Coords) (finish Coords) {
	// fmt.Println("Current start point: ", start)
	visited := map[Coords]bool{
		start: true,
	}

	explosive := explosionArea(gameMap, start)
	// fmt.Println("Bomb Hitting Cells: ", explosive)

	queue := []Coords{start}

	counter := 0
	for len(queue) > 0 {

		// fmt.Println("Queue: ", queue)
		// fmt.Println("Visited Coordinates: ", visited)
		counter++

		coordsDir := []Coords{
			Coords{x: queue[0].x, y: queue[0].y - 1},
			Coords{x: queue[0].x, y: queue[0].y + 1},
			Coords{x: queue[0].x - 1, y: queue[0].y},
			Coords{x: queue[0].x + 1, y: queue[0].y},
		}

		finish = queue[0]

		if !explosive[finish] {
			// fmt.Println("Finishing at: ", finish)
			return finish
		}

		// fmt.Println("Current optimal coord: ", finish)

		queue = queue[1:]

		for _, coord := range coordsDir {
			if (coord.x >= 0 && coord.x <= maxX) && (coord.y >= 0 && coord.y <= maxY) {
				if gameMap[coord.y][coord.x] == '.' && !visited[coord] {
					queue = append(queue, coord)
					visited[coord] = true
				}
			}
		}

	}


	return
}  
