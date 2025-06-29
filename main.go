package main

import (
	"errors"
	"fmt"
	"math"

	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/entities"
	"github.com/pyrdelic/goball/level"
)

type Game struct {
	paddle entities.Paddle
	//bricks []entities.Brick
	ball         entities.Ball
	balls        [config.BallMaxCount]*entities.Ball
	lives        int
	level        *level.Level
	currLevelNum int
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

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// We exit the game by returning a custom error
		return ErrTerminated
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		for _, b := range g.balls {
			if b == nil {
				continue
			}
			b.Grabbed = false
			// make sure the ball launches upwards
			if b.SpeedY > 0 {
				b.SpeedY = -b.SpeedY
			}
		}
		g.ball.Grabbed = false
		if g.ball.SpeedY > 0 {
			// launch upwards
			g.ball.SpeedY = -g.ball.SpeedY
		}
	}

	alreadyBouncedBrick := false // prevents bounce cancellation if multiple collision
	// brick collisions
	for iRow := 0; iRow < config.BrickRowCount; iRow++ {
		for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
			if g.level.Bricks[iRow][iColumn] == nil {
				continue
			}
			if isColliding(&g.ball.Rect, &g.level.Bricks[iRow][iColumn].Rect) {
				// collision detected
				collidedBrick := g.level.Bricks[iRow][iColumn]
				// calculate collision lengths x, y to determine x (L, R) or y (T, B) sided bounce
				if !alreadyBouncedBrick {
					var xCollisionLength, yCollisionLength float64
					// xCollisionLength
					if g.ball.Rect.X < collidedBrick.Rect.X {
						xCollisionLength = g.ball.Rect.X + g.ball.Rect.W - collidedBrick.Rect.X
					} else {
						xCollisionLength = collidedBrick.Rect.X + collidedBrick.Rect.X - g.ball.Rect.X
					}
					// yCollisionLength
					if g.ball.Rect.Y < collidedBrick.Rect.Y {
						yCollisionLength = g.ball.Rect.Y + g.ball.Rect.W - collidedBrick.Rect.Y
					} else {
						yCollisionLength = collidedBrick.Rect.Y + collidedBrick.Rect.H - g.ball.Rect.Y
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
				g.level.TotalHealth--
				collidedBrick.Health--
			}
		}
	}
	// destroy bricks with 0 or less health, decrease level.TotalHealth
	for iRow := 0; iRow < config.BrickRowCount; iRow++ {
		for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
			if g.level.Bricks[iRow][iColumn] == nil {
				continue
			}
			if g.level.Bricks[iRow][iColumn].Health <= 0 {
				g.level.Bricks[iRow][iColumn] = nil
			}
		}
	}
	// g.level.Bricks = slices.DeleteFunc(g.level.Bricks, func(b entities.Brick) bool {
	// 	if b.Health <= 0 {
	// 		return true
	// 	} else {
	// 		return false
	// 	}
	// })

	// wall collisions & bounce
	for _, b := range g.balls {
		if b == nil {
			continue
		}
		// left wall
		if b.Rect.X <= 0 && b.SpeedX < 0 {
			b.SpeedX = -b.SpeedX
		}
		// right wall
		if b.Rect.X+b.Rect.W >= config.PlayAreaWidth && b.SpeedX > 0 {
			b.SpeedX = -b.SpeedX
		}
		// ceiling
		if b.Rect.Y <= 0 && b.SpeedY < 0 {
			b.SpeedY = -b.SpeedY
		}
		// floor
		if b.Rect.Y+b.Rect.H >= config.PlayAreaHeight && b.SpeedY > 0 {
			// TODO: destroy ball
			b.SpeedY = -b.SpeedY
		}
	}

	// left wall
	if g.ball.Rect.X <= 0 && g.ball.SpeedX < 0 {
		g.ball.SpeedX = -g.ball.SpeedX
	}
	// right wall
	if g.ball.Rect.X+g.ball.Rect.W >= config.PlayAreaWidth && g.ball.SpeedX > 0 {
		g.ball.SpeedX = -g.ball.SpeedX
	}
	// ceiling
	if g.ball.Rect.Y <= 0 && g.ball.SpeedY < 0 {
		g.ball.SpeedY = -g.ball.SpeedY
	}
	// floor
	if g.ball.Rect.Y+g.ball.Rect.W >= config.PlayAreaHeight && g.ball.SpeedY > 0 {
		// TODO: lose a ball / lose a life / game over
		g.ball.SpeedY = -g.ball.SpeedY
	}

	// Paddle collisions & bounce
	for _, b := range g.balls {
		if b == nil {
			continue
		}
		if !isColliding(&b.Rect, &g.paddle.Rect) {
			continue
		}
		ballCenterX := b.Rect.X + b.Rect.W/2
		fmt.Println("Ball centerX:", ballCenterX)
		segmentAngleDegrees := 22.5
		paddleSegmentLenX := g.paddle.Rect.W / 6
		fmt.Println("paddleSegmentLenx:", paddleSegmentLenX)

		if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX {
			fmt.Println("multiball hit segment: 1")
			b.CalcXYForAngle(360.0 - segmentAngleDegrees*2 - segmentAngleDegrees/2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*2 {
			fmt.Println("multiball hit segment: 2")
			b.CalcXYForAngle(360.0 - segmentAngleDegrees - segmentAngleDegrees/2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*3 {
			fmt.Println("multiball hit segment: 3")
			b.CalcXYForAngle(360.0 - segmentAngleDegrees/2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*4 {
			fmt.Println("multiball hit segment: 4")
			b.CalcXYForAngle(segmentAngleDegrees / 2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*5 {
			fmt.Println("multiball hit segment: 5")
			b.CalcXYForAngle(segmentAngleDegrees + segmentAngleDegrees/2.0)
		} else {
			fmt.Println("multiball hit segment: 6")
			b.CalcXYForAngle(segmentAngleDegrees*2 + segmentAngleDegrees/2.0)
		}
		fmt.Println("Ball sX, sY:", b.SpeedX, b.SpeedY)
		// ensure that the ball bounces upwards
		if b.SpeedY > 0 {
			b.SpeedY = -b.SpeedY
		}

	}
	// if isColliding(&g.ball.Rect, &g.paddle.Rect) {
	// 	// Ball's bounce angle is determined by the point
	// 	// of collision on the paddle.
	// 	ballCenterX := g.ball.Rect.X + g.ball.Rect.W/2.0

	// 	// segmented determination of bounce (launch) angle
	// 	segmentAngleDegrees := 22.5
	// 	paddleSegmentLenX := g.paddle.Rect.W / 6
	// 	if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX {
	// 		fmt.Println("hit segment 1") // works
	// 		g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(
	// 			g.ball.SpeedBase,
	// 			360.0-(segmentAngleDegrees/2.0)-(segmentAngleDegrees*2))
	// 	} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*2 {
	// 		fmt.Println("hit segment 2") // works
	// 		g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(g.ball.SpeedBase,
	// 			360.0-(segmentAngleDegrees/2.0)-segmentAngleDegrees)
	// 	} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*3 {
	// 		fmt.Println("hit segment 3 ") // launches straight up, why?
	// 		g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(g.ball.SpeedBase,
	// 			360.0-(segmentAngleDegrees/2.0))
	// 		fmt.Println(360.0 - (segmentAngleDegrees / 2.0))
	// 	} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*4 {
	// 		fmt.Println("hit segment 4") // launches straight up, why?
	// 		g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(g.ball.SpeedBase,
	// 			segmentAngleDegrees/2.0)
	// 	} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*5 {
	// 		fmt.Println("hit segment 5") // works
	// 		g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(g.ball.SpeedBase,
	// 			(segmentAngleDegrees/2.0)+segmentAngleDegrees)
	// 	} else {
	// 		fmt.Println("hit segment 6") // works
	// 		g.ball.SpeedX, g.ball.SpeedY = speedXYForAngle(g.ball.SpeedBase,
	// 			(segmentAngleDegrees/2.0)+segmentAngleDegrees*2)
	// 	}
	// 	fmt.Println(segmentAngleDegrees, segmentAngleDegrees/2.0)

	// 	if g.ball.SpeedY > 0 {
	// 		g.ball.SpeedY = -g.ball.SpeedY
	// 	}
	// }

	// for i := range g.bricks {
	// 	g.bricks[i].Update()
	// }

	g.paddle.Update()

	// update balls
	for _, b := range g.balls {
		if b == nil {
			continue
		}
		if b.Grabbed {
			b.Rect.X = g.paddle.Rect.X
		} else {
			b.Update()
		}
	}

	if !g.ball.Grabbed {
		g.ball.Update()
	} else {
		g.ball.Rect.X = g.paddle.Rect.X
	}

	if g.level.TotalHealth <= 0 {
		g.currLevelNum++
		g.level = level.NewLevel(g.currLevelNum)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//ebitenutil.DebugPrint(screen, "Hello, World!")
	for iRow := 0; iRow < config.BrickRowCount; iRow++ {
		for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
			g.level.Bricks[iRow][iColumn].Draw(screen)
		}
	}
	// for i := range g.bricks {
	// 	g.bricks[i].Draw(screen)
	// }
	g.paddle.Draw(screen)
	g.ball.Draw(screen)
	for _, b := range g.balls {
		b.Draw(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"lives: %d\nlvl health: %d",
		g.lives,
		g.level.TotalHealth))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.PlayAreaWidth, config.PlayAreaHeight
}

// Custom error to exit the game loop in a regular way.
var ErrTerminated = errors.New("terminated")

func main() {
	// if true {
	// 	var degree float64 = 22.5
	// 	var baseSpeed float64 = 22.0
	// 	sX, sY := speedXYForAngle(baseSpeed, degree)
	// 	fmt.Println(degree, "|", sX, "|", sY)
	// 	return
	// }
	if false {
		lvl := level.Level{}
		lvl.LoadFromFile("level1.txt")
		return
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GO-BALL")

	game := Game{}

	game.balls[0] = entities.NewBall(
		100.0,
		100.0,
		config.BallStartingSpeed,
		config.BallStartingAngle,
		true,
	)

	game.currLevelNum = 1
	game.level = level.NewLevel(game.currLevelNum)
	game.level.PrintLevel()
	// init paddle
	cursorX, _ := ebiten.CursorPosition()
	game.lives = config.StartingLives
	game.paddle.Rect.X = float64(cursorX)
	game.paddle.Rect.Y = 200
	game.paddle.Rect.W = config.PaddleStartingWidth
	game.paddle.Rect.H = 5
	game.paddle.Image = ebiten.NewImage(int(game.paddle.Rect.W), int(game.paddle.Rect.H))
	game.paddle.Image.Fill(color.White)

	// init ball

	game.ball.Rect.X = 100
	game.ball.Rect.Y = 100
	game.ball.Rect.W = config.BallSize
	game.ball.Rect.H = config.BallSize
	//game.ball.Speed = -2
	game.ball.SpeedBase = config.BallStartingSpeed
	game.ball.SpeedX, game.ball.SpeedY = speedXYForAngle(
		game.ball.SpeedBase, 360.0-12.75)
	game.ball.Grabbed = true

	game.ball.Image = ebiten.NewImage(int(game.ball.Rect.W), int(game.ball.Rect.H))
	game.ball.Image.Fill(color.White)

	ebiten.SetVsyncEnabled(false)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(&game); err != nil {
		if err == ErrTerminated {
			// Regular termination
			return
		}
		log.Fatal(err)
	}
}
