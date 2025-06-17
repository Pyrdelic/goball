package main

import (
	//"crypto/rand"

	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pyrdelic/goball/entities"
)

type Game struct {
	paddleStruct entities.Paddle
	brickStructs []entities.Brick
	ball         entities.Ball
}

func isColliding(a *entities.Ball, b *entities.Brick) bool {
	// x axis
	if !(a.X+a.W < b.X || b.X+b.W < a.X) {

		// y axis
		if !(a.Y+a.H < b.Y || b.Y+b.H < a.Y) {
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
	// check collisions
	// TODO: optimize ( reduce checks per tick)
	if isColliding(&g.ball, &g.brickStructs[0]) {
		fmt.Println("Is colliding")
		g.ball.Speed = -g.ball.Speed
	}

	for i := range g.brickStructs {
		g.brickStructs[i].Update()
	}
	g.paddleStruct.Update()
	g.ball.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	for i := range g.brickStructs {
		g.brickStructs[i].Draw(screen)
	}
	g.paddleStruct.Draw(screen)
	g.ball.Draw(screen)
}

const (
	playAreaWidth    = 320
	playAreaHeight   = 240
	brickColumnCount = 16
	brickRowCount    = 6
	brickCount       = brickColumnCount * brickRowCount
	brickHeight      = 10
	brickWidth       = playAreaWidth / brickColumnCount
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return playAreaWidth, playAreaHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GO-BALL")

	game := Game{}

	// init paddle
	game.paddleStruct.X, _ = ebiten.CursorPosition()
	game.paddleStruct.Y = 150
	game.paddleStruct.W = 80
	game.paddleStruct.H = 10
	game.paddleStruct.Image = ebiten.NewImage(game.paddleStruct.W, game.paddleStruct.H)
	//fmt.Println("dimensions: ", game.paddleStruct.W, game.paddleStruct.H) // works
	game.paddleStruct.Image.Fill(color.White)

	// init ball

	game.ball.X = 200
	game.ball.Y = 0
	game.ball.W = 10
	game.ball.H = 10
	game.ball.Speed = 2
	game.ball.Image = ebiten.NewImage(game.ball.W, game.ball.H)
	game.ball.Image.Fill(color.White)

	// init bricks
	if false {
		for iBrickRow := 0; iBrickRow < brickRowCount; iBrickRow++ {
			for iBrickColumn := 0; iBrickColumn < brickColumnCount; iBrickColumn++ {
				brick := entities.Brick{}
				brick.X = iBrickColumn * brickWidth
				brick.Y = iBrickRow * brickHeight
				brick.W = brickWidth
				brick.H = brickHeight
				brick.Image = ebiten.NewImage(brick.W, brick.H)
				brick.Image.Fill(color.RGBA{
					R: uint8(iBrickRow * 25),
					G: uint8(iBrickColumn * 10),
					B: uint8(127),
					A: uint8(255),
				})
				game.brickStructs = append(game.brickStructs, brick)
			}
		}
	} else { // testing stuff
		brick := entities.Brick{}
		brick.X = 200
		brick.Y = 100
		brick.W = brickWidth
		brick.H = brickHeight
		brick.Image = ebiten.NewImage(brick.W, brick.H)
		brick.Image.Fill(color.RGBA{
			R: uint8(127),
			G: uint8(127),
			B: uint8(127),
			A: uint8(255),
		})
		game.brickStructs = append(game.brickStructs, brick)
	}

	ebiten.SetVsyncEnabled(false)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
