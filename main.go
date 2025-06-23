package main

import (
	//"crypto/rand"
	"fmt"
	"math"

	"image/color"
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pyrdelic/goball/entities"
)

type Game struct {
	paddle entities.Paddle
	bricks []entities.Brick
	ball   entities.Ball
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

// returns a game coordinate -rotated x y speed components
// for a desired angle and base speed multiplier
func speedXYForAngle(speedBase float64, angle float64) (float64, float64) {
	radian := angle * (math.Pi / 180)
	speedX := speedBase * math.Sin(radian)
	speedY := speedBase * math.Cos(radian)
	// flip Y component to correct for game space coordinate system
	fmt.Println(int(speedX), int(-speedY))
	return speedX, -speedY
}

func (g *Game) Update() error {
	// check collisions
	// TODO: optimize ( reduce checks per tick)

	alreadyBouncedBrick := false // prevents bounce cancellation if multiple collision
	// brick collisions
	for i := range g.bricks {
		if isColliding(&g.ball.Rect, &g.bricks[i].Rect) {
			// collision detected

			// calculate collision lengths x, y to determine x (L, R) or y (T, B) sided bounce
			if !alreadyBouncedBrick {
				var xCollisionLength, yCollisionLength float64
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
	if isColliding(&g.ball.Rect, &g.paddle.Rect) {
		// Ball's bounce angle is determined by the point
		// of collision on the paddle.
		ballCenterX := g.ball.Rect.X + g.ball.Rect.W/2

		// segmented determination of bounce (launch) angle
		segmentAngleDegrees := 22.5
		paddleSegmentLenX := g.paddle.Rect.W / 6
		if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX {
			fmt.Println("hit segment 1") // works
			g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(
				g.ball.SpeedMultiplier,
				360.0-(segmentAngleDegrees/2.0)-(segmentAngleDegrees*2))
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*2 {
			fmt.Println("hit segment 2") // works
			g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(g.ball.SpeedMultiplier,
				360.0-(segmentAngleDegrees/2.0)-segmentAngleDegrees)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*3 {
			fmt.Println("hit segment 3 ") // launches straight up, why?
			g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(g.ball.SpeedMultiplier,
				360.0-(segmentAngleDegrees/2.0))
			fmt.Println(360.0 - (segmentAngleDegrees / 2.0))
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*4 {
			fmt.Println("hit segment 4") // launches straight up, why?
			g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(g.ball.SpeedMultiplier,
				segmentAngleDegrees/2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*5 {
			fmt.Println("hit segment 5") // works
			g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(g.ball.SpeedMultiplier,
				(segmentAngleDegrees/2.0)+segmentAngleDegrees)
		} else {
			fmt.Println("hit segment 6") // works
			g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(g.ball.SpeedMultiplier,
				(segmentAngleDegrees/2.0)+segmentAngleDegrees*2)
		}
		fmt.Println(segmentAngleDegrees, segmentAngleDegrees/2.0)

		if g.ball.SpeedY > 0 {
			g.ball.SpeedY = -g.ball.SpeedY
		}
	}

	for i := range g.bricks {
		g.bricks[i].Update()
	}
	g.paddle.Update()
	g.ball.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	for i := range g.bricks {
		g.bricks[i].Draw(screen)
	}
	g.paddle.Draw(screen)
	g.ball.Draw(screen)
}

const (
	playAreaHeight   = 240 // in-game resolution
	playAreaWidth    = 320 // in-game resolution
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
	// if true {
	// 	var degree float64 = 22.5
	// 	var baseSpeed float64 = 22.0
	// 	sX, sY := speedXYForAngle(baseSpeed, degree)
	// 	fmt.Println(degree, "|", sX, "|", sY)
	// 	return
	// }

	fmt.Println("brick width: ", brickWidth)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GO-BALL")

	game := Game{}

	// init paddle
	cursorX, _ := ebiten.CursorPosition()
	game.paddle.Rect.X = float64(cursorX)
	game.paddle.Rect.Y = 200
	game.paddle.Rect.W = 80
	game.paddle.Rect.H = 10
	game.paddle.Image = ebiten.NewImage(int(game.paddle.Rect.W), int(game.paddle.Rect.H))
	game.paddle.Image.Fill(color.White)

	// init ball

	game.ball.Rect.X = 100
	game.ball.Rect.Y = 100
	game.ball.Rect.W = 10
	game.ball.Rect.H = 10
	//game.ball.Speed = -2
	game.ball.SpeedX = -2
	game.ball.SpeedY = -2
	game.ball.SpeedMultiplier = 6.0
	game.ball.Image = ebiten.NewImage(int(game.ball.Rect.W), int(game.ball.Rect.H))
	game.ball.Image.Fill(color.White)

	// init bricks
	if true {
		for iBrickRow := 0; iBrickRow < brickRowCount; iBrickRow++ {
			for iBrickColumn := 0; iBrickColumn < brickColumnCount; iBrickColumn++ {
				brick := entities.Brick{}
				brick.Health = 1
				brick.Rect.X = float64(iBrickColumn * brickWidth)
				brick.Rect.Y = float64(iBrickRow * brickHeight)
				brick.Rect.W = brickWidth
				brick.Rect.H = brickHeight
				brick.Image = ebiten.NewImage(int(brick.Rect.W), int(brick.Rect.H))
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
		brick.Image = ebiten.NewImage(int(brick.Rect.W), int(brick.Rect.H))
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
