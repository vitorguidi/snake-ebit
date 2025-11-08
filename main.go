package main

import (
	"image/color"
	"log"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	squareSide  = 20
	pixelWidth  = 800
	pixelHeight = 600
	interval    = 150 * time.Millisecond
	gridWidth   = pixelWidth / squareSide
	gridHeight  = pixelHeight / squareSide
)

var (
	directions = map[ebiten.Key]point{
		ebiten.KeyW: {0, -1},
		ebiten.KeyA: {-1, 0},
		ebiten.KeyS: {0, 1},
		ebiten.KeyD: {1, 0},
	}
)

type point struct {
	x, y int
}

type snake struct {
	body      []point
	direction point
}

func randomPoint() point {
	return point{rand.IntN(gridWidth), rand.IntN(gridHeight)}
}

func newSnake() *snake {
	return &snake{
		body:      []point{randomPoint()},
		direction: point{0, 0},
	}
}

func (s *snake) head() point {
	return s.body[len(s.body)-1]
}

func (s *snake) changeDirection(direction point) {
	s.direction = direction
}

// returns true if ate food, false otherwise
func (s *snake) advance(foodPos point) bool {
	if s.direction.x == 0 && s.direction.y == 0 {
		return false
	}
	head := s.head()
	newX := (head.x + s.direction.x + gridWidth) % gridWidth
	newY := (head.y + s.direction.y + gridHeight) % gridHeight
	s.body = append(s.body, point{newX, newY})
	if head != foodPos {
		s.body = s.body[1:]
		return false
	}
	return true
}

func (s *snake) collision() bool {
	return slices.Contains(s.body[:len(s.body)-1], s.head())
}

// Game implements ebiten.Game interface.
type Game struct {
	snake          snake
	food           point
	updateInterval time.Duration
}

func (g *Game) changeSnakeDirection() {
	keys := []ebiten.Key{ebiten.KeyA, ebiten.KeyW, ebiten.KeyS, ebiten.KeyD}
	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			g.snake.changeDirection(directions[key])
			return
		}
	}
}

func (g *Game) spawnFood() {
	for {
		pos := randomPoint()
		if slices.Contains(g.snake.body, pos) {
			continue
		}
		g.food = pos
		break
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.changeSnakeDirection()
	time.Sleep(g.updateInterval)
	if g.snake.advance(g.food) {
		g.spawnFood()
	}
	if g.snake.collision() {
		g.snake = *newSnake()
		g.spawnFood()
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.snake.body {
		vector.DrawFilledRect(
			screen,
			float32(p.x*squareSide),
			float32(p.y*squareSide),
			squareSide,
			squareSide,
			color.White,
			true,
		)
	}
	vector.DrawFilledRect(
		screen,
		float32(g.food.x*squareSide),
		float32(g.food.y*squareSide),
		squareSide,
		squareSide,
		color.White,
		true,
	)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return pixelWidth, pixelHeight
}

func main() {
	game := &Game{
		snake:          *newSnake(),
		updateInterval: interval,
	}
	game.spawnFood()
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(pixelWidth, pixelHeight)
	ebiten.SetWindowTitle("La cobrita")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
