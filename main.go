package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 800
	tileSize     = 20
	gridTop      = 20
	gridWidth    = screenWidth / tileSize
	gridHeight   = screenHeight / tileSize
	gridSize     = screenWidth / tileSize
)

var (
	yellow = color.RGBA{255, 255, 0, 255}
	grey   = color.RGBA{128, 128, 128, 255}
	black  = color.RGBA{0, 0, 0, 255}
	//black  = [4]float64{0, 0, 0, 1}
)

type World struct {
	screenWidth  int
	screenHeight int
	tileSize     int
	gridWidth    int
	gridHeight   int
	gridSize     int
	gridTop      int
	alive        bool
	liveCells    map[tile]struct{}
	isSimulating bool
	lastUpdate   time.Time
}

type tile struct {
	x, y int
}

// NewWorld creates a new world
func NewWorld(screenWidth, screenHeight, tileSize int) *World {
	return &World{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		tileSize:     tileSize,
		gridWidth:    gridWidth,
		gridHeight:   gridHeight,
		gridSize:     gridSize,
		gridTop:      gridTop,
		liveCells:    make(map[tile]struct{}),
		isSimulating: false,
		alive:        false,
		lastUpdate:   time.Now(),
	}
}

// DrawWorld draws the world
func (w *World) DrawWorld(screen *ebiten.Image) {

	// Draw the lines of the grid
	for i := 0; i <= w.gridSize; i++ {
		thickness := float32(1.0)

		// Vertical lines
		x := float32(i * w.tileSize)
		vector.StrokeLine(
			screen,
			x,
			float32(0),
			x,
			float32(w.gridTop+(w.gridSize*w.tileSize)), // Fix grid height calculation
			thickness,
			black,
			false,
		)

		// Horizontal lines
		y := float32(w.gridTop + i*w.tileSize)
		vector.StrokeLine(
			screen,
			0,
			y,
			float32(w.gridSize*w.tileSize), // Fix grid width calculation
			y,
			thickness,
			black,
			false,
		)
	}

}

func (w *World) handleMouseClick(x, y int) {
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return
	}

	// Calculate the cell clicked
	cellX := x / w.tileSize
	cellY := (y - w.gridTop) / w.tileSize

	if cellX < 0 || cellX >= w.gridWidth || cellY < 0 || cellY >= w.gridHeight {
		return
	}
	clickedCell := tile{x: cellX, y: cellY}
	if _, isAlive := w.liveCells[clickedCell]; isAlive {
		delete(w.liveCells, clickedCell)
	} else {
		w.liveCells[clickedCell] = struct{}{}
	}

}

// fillCell draws a cell filled with a color
func (w *World) fillCell(screen *ebiten.Image, x, y int, color color.Color) {
	vector.DrawFilledRect(screen, float32(x*w.tileSize), float32(w.gridTop+y*w.tileSize), float32(w.tileSize), float32(w.tileSize), yellow, false)
}

// drawliveCells draws all the live cells
func (w *World) drawLiveCells(screen *ebiten.Image) {
	for cell := range w.liveCells {
		w.fillCell(screen, cell.x, cell.y, yellow)
	}

}

// generateRandomCells generates random cells
func (w *World) generateRandomCells() {
	// Clear the current cells
	w.liveCells = make(map[tile]struct{})

	//time as seed
	rand.New(rand.NewSource(time.Now().UnixNano()))
	totalCells := w.gridWidth * w.gridHeight
	numCells := rand.Intn((totalCells / 5) + totalCells/5)

	for i := 0; i < numCells; i++ {
		x := rand.Intn(w.gridWidth)
		y := rand.Intn(w.gridHeight)
		w.liveCells[tile{x: x, y: y}] = struct{}{}
	}

}

// SimulateWorld simulates the world following the rules of the game of life.
func (w *World) SimulateWorld() {
	// Create a new map to store the next generation of cells
	nextGeneration := make(map[tile]struct{})
	// Iterate over all the cells
	for cell := range w.liveCells {
		// Count the number of live neighbors
		liveNeighbors := w.countLiveNeighbors(cell.x, cell.y)
		// If the cell has 2 or 3 live neighbors, it survives
		if liveNeighbors == 2 || liveNeighbors == 3 {
			nextGeneration[cell] = struct{}{}
		}
		// Check the neighbors of the cell
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				// Skip the cell itself
				if i == 0 && j == 0 {
					continue
				}
				// Calculate the coordinates of the neighbor
				neighborX := cell.x + i
				neighborY := cell.y + j
				// Count the number of live neighbors
				liveNeighbors := w.countLiveNeighbors(neighborX, neighborY)
				// If the neighbor has exactly 3 live neighbors, it becomes alive
				if liveNeighbors == 3 {
					nextGeneration[tile{x: neighborX, y: neighborY}] = struct{}{}
				}
			}
		}
	}
	// Update the live cells
	w.isSimulating = true
	w.liveCells = nextGeneration
}

// countLiveNeighbors counts the number of live neighbors of a cell
func (w *World) countLiveNeighbors(x, y int) int {
	// Initialize the counter
	liveNeighbors := 0
	// Check the neighbors of the cell
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Skip the cell itself
			if i == 0 && j == 0 {
				continue
			}
			// Calculate the coordinates of the neighbor
			neighborX := x + i
			neighborY := y + j
			// Check if the neighbor is alive
			if _, isAlive := w.liveCells[tile{x: neighborX, y: neighborY}]; isAlive {
				liveNeighbors++
			}
		}
	}
	// Return the number of live neighbors
	return liveNeighbors
}

// generateGosperGliderGun generates a Gosper Glider Gun
func (w *World) generateGosperGliderGun() {
	// Clear the current cells
	w.liveCells = make(map[tile]struct{})
	// Gosper Glider Gun
	gliderGun := []tile{
		{1, 5}, {1, 6}, {2, 5}, {2, 6},
		{11, 5}, {11, 6}, {11, 7},
		{12, 4}, {12, 8},
		{13, 3}, {13, 9},
		{14, 3}, {14, 9},
		{15, 6},
		{16, 4}, {16, 8},
		{17, 5}, {17, 6}, {17, 7},
		{18, 6},
		{21, 3}, {21, 4}, {21, 5},
		{22, 3}, {22, 4}, {22, 5},
		{23, 2}, {23, 6},
		{25, 1}, {25, 2}, {25, 6}, {25, 7},
		{35, 3}, {35, 4},
		{36, 3}, {36, 4},
	}
	for _, cell := range gliderGun {
		w.liveCells[cell] = struct{}{}
	}
}

type Game struct {
	world *World
}

func (g *Game) Update() error {
	// exit game on escape or q key
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	// handle start on g key. generate random cells
	if ebiten.IsKeyPressed(ebiten.KeyG) {
		g.world.generateRandomCells()

	}
	// handle reset on r key
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.world.liveCells = make(map[tile]struct{})
		g.world.isSimulating = false
	}

	// handle space key or s to start simulation
	if ebiten.IsKeyPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.world.SimulateWorld()
	}

	// handle pause on p key
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		g.world.isSimulating = false
	}

	// handle glider gun on 1 key
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.world.generateGosperGliderGun()
	}

	// Run the simulation every 200ms if the simulation is running
	if g.world.isSimulating && time.Since(g.world.lastUpdate) > 300*time.Millisecond {
		g.world.SimulateWorld()
		g.world.lastUpdate = time.Now()
	}

	// handle mouse click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.world.handleMouseClick(x, y)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(grey)
	g.world.DrawWorld(screen)
	g.world.drawLiveCells(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 800
}

func main() {
	// Initialize the world
	world := NewWorld(screenWidth, screenHeight, tileSize)
	game := &Game{world: world}
	ebiten.SetWindowSize(840, 840)
	ebiten.SetWindowTitle("Game Of Life!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

//
