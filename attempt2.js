const maxX = 12,
  maxY = 10,
  bombRadius = 2,
  boxCost = 11,
  step = 1;

const oppositeDirs = {
  left: "right",
  right: "left",
  up: "down",
  down: "up",
};

const printMap = `..........;;;
.!;!;!;!;!;!.
;..........;!
;!;!;!;!;!!!.
;..;.....;..;
.!;!;!;!;!;!;
;....;.;.;!.;
;!.!;!;!;!.!.
;.;.;....;!..
.!!!.!;!;!.!;
.;...;;;...;;`
  .split("\n")
  .map((item) => item.split(""));

const map = `..........;;;
.!;!;!;!;!;!.
;..........;!
;!;!;!;!;!!!.
;..;.....;..;
.!;!;!;!;!;!;
;....;.;.;!.;
;!.!;!;!;!.!.
;.;.;....;!..
.!!!.!;!;!.!;
.;...;;;...;;`
  .split("\n")
  .map((item) => item.split(""));

const getDistance = (current, target) =>
  Math.abs(target.x - current.x) + Math.abs(target.y - current.y);

const Tile = ({
  x,
  y,
  cost = 0,
  target,
  parent = null,
  direction,
  isBox = false,
}) => {
  let distance = 0;

  if (target) distance = getDistance({ x: x, y: y }, target);

  return {
    x: x,
    y: y,
    cost: cost,
    distance: distance,
    costDistance: 0,
    parent: parent,
    direction,
    isBox,
  };
};

const getNeighbours = (current, target) => {
  return [
    Tile({
      x: current.x,
      y: current.y - step,
      cost: current.cost + step,
      target: target,
      direction: "up",
    }),
    Tile({
      x: current.x,
      y: current.y + step,
      cost: current.cost + step,
      target: target,
      direction: "down",
    }),
    Tile({
      x: current.x - step,
      y: current.y,
      cost: current.cost + step,
      target: target,
      direction: "left",
    }),
    Tile({
      x: current.x + step,
      y: current.y,
      cost: current.cost + step,
      target: target,
      direction: "right",
    }),
  ];
};

const contains = (arr, coord) =>
  arr.some((item) => item.x === coord.x && item.y === coord.y);

const findSafeCoord = (start) => {
  const explosiveCells = (bombCoords) => {
    let explosives = [bombCoords];

    // Check Left
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let explCell = {
        x: bombCoords.x - bCounter,
        y: bombCoords.y,
      };

      if (0 <= explCell.x && map[explCell.y][explCell.x] === ".") {
        explosives.push(explCell);
      } else {
        break;
      }
    }

    // Check Right
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let explCell = {
        x: bombCoords.x + bCounter,
        y: bombCoords.y,
      };

      if (explCell.x <= maxX && map[explCell.y][explCell.x] === ".") {
        explosives.push(explCell);
      } else {
        break;
      }
    }

    // Check Up
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let explCell = {
        x: bombCoords.x,
        y: bombCoords.y - bCounter,
      };

      if (0 <= explCell.y && map[explCell.y][explCell.x] === ".") {
        explosives.push(explCell);
      } else {
        break;
      }
    }

    // Check Down
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let explCell = {
        x: bombCoords.x,
        y: bombCoords.y + bCounter,
      };

      if (explCell.y <= maxY && map[explCell.y][explCell.x] === ".") {
        explosives.push(explCell);
      } else {
        break;
      }
    }
    return explosives;
  };

  const explosives = explosiveCells(start);

  const bfs = (start, checker) => {
    let visited = [start];
    let queue = [start];
    let finish;

    while (queue.length > 0) {
      let current = queue.shift();

      const neighbours = getNeighbours(current);

      finish = current;

      if (checker(finish)) {
        return finish;
      }

      neighbours.forEach((coord) => {
        if (
          coord.x >= 0 &&
          coord.x <= maxX &&
          coord.y >= 0 &&
          coord.y <= maxY
        ) {
          if (map[coord.y][coord.x] === "." && !contains(visited, coord)) {
            queue.push(coord);
            visited.push(coord);
          }
        }
      });
    }
    console.error("Failed");
    return;
  };

  return bfs(start, (finish) => !contains(explosives, finish));
};

const findClosestBoxes = (start) => {
  const closestBoxes = (current) => {
    let boxes = [];

    // Check Left
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let boxCell = {
        x: current.x - bCounter,
        y: current.y,
      };

      if (0 <= boxCell.x) {
        if (map[boxCell.y][boxCell.x] === ";") {
          boxes.push(boxCell);
          break;
        } else if (map[boxCell.y][boxCell.x] === "!") {
          break;
        }
      }
    }

    // Check Right
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let boxCell = {
        x: current.x + bCounter,
        y: current.y,
      };

      if (boxCell.x <= maxX) {
        if (map[boxCell.y][boxCell.x] === ";") {
          boxes.push(boxCell);
          break;
        } else if (map[boxCell.y][boxCell.x] === "!") {
          break;
        }
      }
    }

    // Check Up
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let boxCell = {
        x: current.x,
        y: current.y - bCounter,
      };

      if (0 <= boxCell.y) {
        if (map[boxCell.y][boxCell.x] === ";") {
          boxes.push(boxCell);
          break;
        } else if (map[boxCell.y][boxCell.x] === "!") {
          break;
        }
      }
    }

    // Check Down
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let boxCell = {
        x: current.x,
        y: current.y + bCounter,
      };

      if (boxCell.y <= maxY) {
        if (map[boxCell.y][boxCell.x] === ";") {
          boxes.push(boxCell);
          break;
        } else if (map[boxCell.y][boxCell.x] === "!") {
          break;
        }
      }
    }
    return boxes;
  };

  let finish, finBoxes;

  const bfs = (start) => {
    let visited = [start];
    let queue = [start];
    let len = 0;

    while (queue.length > 0) {
      let current = queue.shift();

      const neighbours = getNeighbours(current);
      const boxes = closestBoxes(current);

      if (boxes.length > len) {
        finBoxes = boxes;
        len = boxes.length;
        finish = current;
      }

      neighbours.forEach((coord) => {
        if (
          coord.x >= 0 &&
          coord.x <= maxX &&
          coord.y >= 0 &&
          coord.y <= maxY
        ) {
          if (map[coord.y][coord.x] === "." && !contains(visited, coord)) {
            queue.push(coord);
            visited.push(coord);
          }
        }
      });
    }
    return {
      coord: finish,
      boxes: finBoxes,
    };
  };

  return bfs(start);
};

const buildPath = (tile) => {
  let path = [];
  while (tile.parent) {
    printMap[tile.y][tile.x] = "*";
    if (tile.direction) path.unshift(tile.direction);
    tile = tile.parent;
  }

  // Retreat path building case
  // path.unshift("bomb");
  // path.push(path.length === bombRadius + 1 ? "stay" : "stay", "stay");
  return path;
};

const astar = (start, finish, retreat) => {
  let visitedTiles = [];
  let activeTiles = [start];

  while (activeTiles.length > 0) {
    activeTiles.sort((a, b) => a.costDistance - b.costDistance);

    let current = activeTiles.shift();

    if (current.x === finish.x && current.y === finish.y) {
      console.log("Path found");
      return buildPath(current, retreat);
    }

    let neighbours = getNeighbours(current, finish).filter((nTile) => {
      if (
        nTile.x >= 0 &&
        nTile.x <= maxX &&
        nTile.y >= 0 &&
        nTile.y <= maxY &&
        map[nTile.y][nTile.x] !== "!"
      ) {
        return true;
      }
      return false;
    });

    visitedTiles.push(current);

    for (let i = 0; i < neighbours.length; i++) {
      let nTile = neighbours[i];

      if (map[nTile.y][nTile.x] === ";") nTile.isBox = true;

      if (
        visitedTiles.some((vTile) => vTile.x === nTile.x && vTile.y === nTile.y)
      ) {
        continue;
      }

      if (map[nTile.y][nTile.x] === ";") {
        nTile.cost += boxCost;
      }

      nTile.costDistance = nTile.cost + nTile.distance;
      nTile.parent = current;

      activeTiles.push(nTile);
    }
  }
  console.log("Failed");
  return;
};

const main = () => {
  let finish = Tile({ x: 12, y: 10 });
  let start = Tile({ x: 0, y: 0, target: finish });
  // console.log(findClosestBox(Tile({ x: 0, y: 0 })));
  // return;

  setInterval(() => {
    console.clear();
    // Find destination
    finish = findClosestBoxes(start);
    // return

    // Build path
    let path = astar(start, finish.coord);

    // Find bomb safe cell
    let safeCoord = findSafeCoord(finish.coord);

    // Build retreat path
    let retreatPath = astar(finish.coord, safeCoord);

    // Set bomb
    path.push("bomb");
    path = [...path, ...retreatPath];
    path.push("stay", "stay");
    if (retreatPath.length <= bombRadius) path.push("stay");

    // Update map
    finish.boxes.forEach((box) => {
      map[box.y][box.x] = ".";
      printMap[box.y][box.x] = ".";
    });

    // Update start point
    start = Tile({ x: safeCoord.x, y: safeCoord.y, target: finish });

    // Print map
    printMap.forEach((col) => console.log(col.join("")));

    // console.log(start)
    // Do this again

    // console.log(path);
  }, 2000);
};

main();
