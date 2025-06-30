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
	//ball         entities.Ball
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
	}

	alreadyBouncedBrick := false // prevents bounce cancellation if multiple collision
	// brick collisions

	for _, b := range g.balls {
		if b == nil {
			continue
		}
		for iRow := 0; iRow < config.BrickRowCount; iRow++ {
			for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
				if g.level.Bricks[iRow][iColumn] == nil {
					continue
				}
				if isColliding(&b.Rect, &g.level.Bricks[iRow][iColumn].Rect) {
					collidedBrick := g.level.Bricks[iRow][iColumn]
					// bounce if not already bounced (prevents bounce cancellation)
					if !alreadyBouncedBrick {
						// calculate collision lengts of x and y,
						// this determines if the collision is x or y sided
						// x
						var xCollisionLength, yCollisionLength float64
						if b.Rect.X < collidedBrick.Rect.X {
							xCollisionLength = b.Rect.X + b.Rect.W - collidedBrick.Rect.X
						} else {
							xCollisionLength = collidedBrick.Rect.X + collidedBrick.Rect.X - b.Rect.X
						}
						// y
						if b.Rect.Y < collidedBrick.Rect.Y {
							yCollisionLength = b.Rect.Y + b.Rect.H - collidedBrick.Rect.Y
						} else {
							yCollisionLength = collidedBrick.Rect.Y + collidedBrick.Rect.H - b.Rect.Y
						}

						if xCollisionLength >= yCollisionLength {
							// y-sided collision
							b.SpeedY = -b.SpeedY
							alreadyBouncedBrick = true
						} else {
							// x-sided collision
							b.SpeedX = -b.SpeedX
							alreadyBouncedBrick = true
						}

					}
					collidedBrick.Health--
					g.level.TotalHealth--
				}
			}
		}
	}

	// destroy bricks with 0 or less health
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
			//fmt.Println("multiball hit segment: 1")
			b.CalcXYForAngle(360.0 - segmentAngleDegrees*2 - segmentAngleDegrees/2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*2 {
			//fmt.Println("multiball hit segment: 2")
			b.CalcXYForAngle(360.0 - segmentAngleDegrees - segmentAngleDegrees/2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*3 {
			// fmt.Println("multiball hit segment: 3")
			b.CalcXYForAngle(360.0 - segmentAngleDegrees/2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*4 {
			// fmt.Println("multiball hit segment: 4")
			b.CalcXYForAngle(segmentAngleDegrees / 2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*5 {
			// fmt.Println("multiball hit segment: 5")
			b.CalcXYForAngle(segmentAngleDegrees + segmentAngleDegrees/2.0)
		} else {
			// fmt.Println("multiball hit segment: 6")
			b.CalcXYForAngle(segmentAngleDegrees*2 + segmentAngleDegrees/2.0)
		}
		// ensure that the ball bounces upwards
		if b.SpeedY > 0 {
			b.SpeedY = -b.SpeedY
		}

	}

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

	if g.level.TotalHealth <= 0 {
		g.currLevelNum++
		g.level = level.NewLevel(g.currLevelNum)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for iRow := 0; iRow < config.BrickRowCount; iRow++ {
		for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
			g.level.Bricks[iRow][iColumn].Draw(screen)
		}
	}

	g.paddle.Draw(screen)

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

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GO-BALL")

	game := Game{}

	// init balls
	game.balls[0] = entities.NewBall(
		100.0,
		100.0,
		config.BallStartingSpeed,
		config.BallStartingAngle,
		true,
	)

	// init paddle
	cursorX, _ := ebiten.CursorPosition()
	game.lives = config.StartingLives
	game.paddle.Rect.X = float64(cursorX)
	game.paddle.Rect.Y = 200
	game.paddle.Rect.W = config.PaddleStartingWidth
	game.paddle.Rect.H = 5
	game.paddle.Image = ebiten.NewImage(int(game.paddle.Rect.W), int(game.paddle.Rect.H))
	game.paddle.Image.Fill(color.White)

	// init level
	game.currLevelNum = 1
	game.level = level.NewLevel(game.currLevelNum)
	game.level.PrintLevel()

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
