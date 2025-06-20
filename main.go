package main

import (
	//"crypto/rand"

	"image/color"
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pyrdelic/goball/entities"
)

type Game struct {
	paddleStruct entities.Paddle
	bricks       []entities.Brick
	ball         entities.Ball
}

// Detects a general collision between two Rects
func isColliding(a *entities.Rect, b *entities.Rect) bool {
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

	alreadyBouncedBrick := false // prevents bounce cancellation if multiple collision
	// brick collisions
	for i := range g.bricks {
		if isColliding(&g.ball.Rect, &g.bricks[i].Rect) {
			// collision detected
			// calculate collision lengths x, y to determine x or y sided bounce
			if !alreadyBouncedBrick {
				var xCollisionLength, yCollisionLength int
				// xCollisionLength
				if g.ball.Rect.X < g.bricks[i].Rect.X {
					xCollisionLength = g.ball.Rect.X + g.ball.Rect.W - g.bricks[i].Rect.X
				} else {
					xCollisionLength = g.bricks[i].Rect.X + g.bricks[i].Rect.X - g.ball.Rect.X
				}
				// yCollisionLength
				if g.ball.Rect.Y < g.bricks[i].Rect.Y {
					yCollisionLength = g.ball.Rect.Y + g.ball.Rect.W - g.bricks[i].Rect.Y
				} else {
					yCollisionLength = g.bricks[i].Rect.Y + g.bricks[i].Rect.H - g.ball.Rect.Y
				}

				if xCollisionLength >= yCollisionLength {
					// top / bottom collision
					g.ball.SpeedY = -g.ball.SpeedY
					alreadyBouncedBrick = true
				} else {
					// side collision
					g.ball.SpeedX = -g.ball.SpeedX
					alreadyBouncedBrick = true
				}
			}
			g.bricks[i].Health--
		}
	}
	// destroy bricks with 0 or less health
	g.bricks = slices.DeleteFunc(g.bricks, func(b entities.Brick) bool {
		if b.Health <= 0 {
			return true
		} else {
			return false
		}
	})

	// wall collisions & bounce
	// left wall
	if g.ball.Rect.X <= 0 {
		g.ball.SpeedX = -g.ball.SpeedX
	}
	// right wall
	if g.ball.Rect.X+g.ball.Rect.W >= playAreaWidth {
		g.ball.SpeedX = -g.ball.SpeedX
	}
	// ceiling
	if g.ball.Rect.Y <= 0 {
		g.ball.SpeedY = -g.ball.SpeedY
	}
	// floor
	if g.ball.Rect.Y+g.ball.Rect.W >= playAreaHeight {
		// TODO: game over / losing a ball
		g.ball.SpeedY = -g.ball.SpeedY
	}

	// Paddle collisions & bounce
	if isColliding(&g.ball.Rect, &g.paddleStruct.Rect) {
		if g.ball.SpeedY > 0 {
			g.ball.SpeedY = -g.ball.SpeedY
		}
	}

	for i := range g.bricks {
		g.bricks[i].Update()
	}
	g.paddleStruct.Update()
	g.ball.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	for i := range g.bricks {
		g.bricks[i].Draw(screen)
	}
	g.paddleStruct.Draw(screen)
	g.ball.Draw(screen)
}

const (
	playAreaWidth    = 320 // in-game resolution
	playAreaHeight   = 240 // in-game resolution
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
	game.paddleStruct.Rect.X, _ = ebiten.CursorPosition()
	game.paddleStruct.Rect.Y = 150
	game.paddleStruct.Rect.W = 80
	game.paddleStruct.Rect.H = 10
	game.paddleStruct.Image = ebiten.NewImage(game.paddleStruct.Rect.W, game.paddleStruct.Rect.H)
	game.paddleStruct.Image.Fill(color.White)

	// init ball

	game.ball.Rect.X = 150
	game.ball.Rect.Y = 150
	game.ball.Rect.W = 10
	game.ball.Rect.H = 10
	//game.ball.Speed = -2
	game.ball.SpeedX = -2
	game.ball.SpeedY = -2
	game.ball.Image = ebiten.NewImage(game.ball.Rect.W, game.ball.Rect.H)
	game.ball.Image.Fill(color.White)

	// init bricks
	if true {
		for iBrickRow := 0; iBrickRow < brickRowCount; iBrickRow++ {
			for iBrickColumn := 0; iBrickColumn < brickColumnCount; iBrickColumn++ {
				brick := entities.Brick{}
				brick.Health = 1
				brick.Rect.X = iBrickColumn * brickWidth
				brick.Rect.Y = iBrickRow * brickHeight
				brick.Rect.W = brickWidth
				brick.Rect.H = brickHeight
				brick.Image = ebiten.NewImage(brick.Rect.W, brick.Rect.H)
				brick.Image.Fill(color.RGBA{
					R: uint8(iBrickRow * 25),
					G: uint8(iBrickColumn * 10),
					B: uint8(127),
					A: uint8(255),
				})
				game.bricks = append(game.bricks, brick)
			}
		}
	} else { // testing stuff
		brick := entities.Brick{}
		brick.Rect.X = 200
		brick.Rect.Y = 100
		brick.Rect.W = brickWidth
		brick.Rect.H = brickHeight
		brick.Image = ebiten.NewImage(brick.Rect.W, brick.Rect.H)
		brick.Image.Fill(color.RGBA{
			R: uint8(127),
			G: uint8(127),
			B: uint8(127),
			A: uint8(255),
		})
		game.bricks = append(game.bricks, brick)
	}

	ebiten.SetVsyncEnabled(false)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
