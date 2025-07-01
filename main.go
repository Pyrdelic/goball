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
	paddle *entities.Paddle
	//bricks []entities.Brick
	//ball         entities.Ball
	balls        [config.BallMaxCount]*entities.Ball
	BallCount    int
	lives        int
	level        *level.Level
	currLevelNum int
	GameOver     bool
}

func UpdateNode(n entities.Node) {
	n.Update()
}

func DrawNode(n entities.Node) {

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
	//fmt.Println(int(speedX), int(-speedY))
	return speedX, -speedY
}

func (g *Game) Update() error {
	// check collisions
	// TODO: optimize ( reduce checks per tick)

	// // test
	// fmt.Println(g.balls)
	// asdf := 0
	// if true {
	// 	for i := 0; i < len(g.balls); i++ {
	// 		g.balls[asdf] = nil
	// 	}

	// }
	// fmt.Println(g.balls)

	if g.BallCount <= 0 {
		// lose a life
		g.lives--
		if g.lives >= 0 {
			// initialize a new ball
			g.InitGrabbedBall()
			g.BallCount++
		} else {
			fmt.Println("Out of lives.")
			fmt.Println("GAME OVER")
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// We exit the game by returning a custom error
		return ErrTerminated
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		for _, b := range g.balls {
			if b == nil {
				continue
			}
			if b.Grabbed {
				b.Grabbed = false
				// make sure the ball launches upwards
				if b.SpeedY > 0 {
					b.SpeedY = -b.SpeedY
				}
			}
		}
	}

	alreadyBouncedBrick := false // prevents bounce cancellation if multiple collision
	// brick collisions

	for i := 0; i < len(g.balls); i++ {
		if g.balls[i] == nil {
			continue
		}
		for iRow := 0; iRow < config.BrickRowCount; iRow++ {
			for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
				if g.level.Bricks[iRow][iColumn] == nil {
					continue
				}
				if isColliding(&g.balls[i].Rect, &g.level.Bricks[iRow][iColumn].Rect) {
					collidedBrick := g.level.Bricks[iRow][iColumn]
					// bounce if not already bounced (prevents bounce cancellation)
					if !alreadyBouncedBrick {
						// calculate collision lengts of x and y,
						// this determines if the collision is x or y sided
						// x
						var xCollisionLength, yCollisionLength float64
						if g.balls[i].Rect.X < collidedBrick.Rect.X {
							xCollisionLength = g.balls[i].Rect.X + g.balls[i].Rect.W - collidedBrick.Rect.X
						} else {
							xCollisionLength = collidedBrick.Rect.X + collidedBrick.Rect.X - g.balls[i].Rect.X
						}
						// y
						if g.balls[i].Rect.Y < collidedBrick.Rect.Y {
							yCollisionLength = g.balls[i].Rect.Y + g.balls[i].Rect.H - collidedBrick.Rect.Y
						} else {
							yCollisionLength = collidedBrick.Rect.Y + collidedBrick.Rect.H - g.balls[i].Rect.Y
						}

						if xCollisionLength >= yCollisionLength {
							// y-sided collision
							g.balls[i].SpeedY = -g.balls[i].SpeedY
							alreadyBouncedBrick = true
						} else {
							// x-sided collision
							g.balls[i].SpeedX = -g.balls[i].SpeedX
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
	for i := 0; i < len(g.balls); i++ {
		if g.balls[i] == nil {
			continue
		}
		// left wall
		if g.balls[i].Rect.X <= 0 && g.balls[i].SpeedX < 0 {
			g.balls[i].SpeedX = -g.balls[i].SpeedX
		}
		// right wall
		if g.balls[i].Rect.X+g.balls[i].Rect.W >= config.PlayAreaWidth &&
			g.balls[i].SpeedX > 0 {
			g.balls[i].SpeedX = -g.balls[i].SpeedX
		}
		// ceiling
		if g.balls[i].Rect.Y <= 0 && g.balls[i].SpeedY < 0 {
			g.balls[i].SpeedY = -g.balls[i].SpeedY
		}
		// floor
		if g.balls[i].Rect.Y+g.balls[i].Rect.H >= config.PlayAreaHeight && g.balls[i].SpeedY > 0 {
			// TODO: destroy ball
			if config.GodMode {
				// god mode just bounces off the floor too
				g.balls[i].SpeedY = -g.balls[i].SpeedY
			} else {
				g.balls[i] = nil
				//fmt.Println("Ball destroyed")
				g.BallCount--
			}
		}
	}

	// Paddle collisions & bounce
	for i := 0; i < len(g.balls); i++ {
		if g.balls[i] == nil {
			continue
		}
		if !isColliding(&g.balls[i].Rect, &g.paddle.Rect) {
			continue
		}
		ballCenterX := g.balls[i].Rect.X + g.balls[i].Rect.W/2
		//fmt.Println("Ball centerX:", ballCenterX)
		segmentAngleDegrees := 22.5
		paddleSegmentLenX := g.paddle.Rect.W / 6
		//fmt.Println("paddleSegmentLenx:", paddleSegmentLenX)

		if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX {
			//fmt.Println("multiball hit segment: 1")
			g.balls[i].CalcXYForAngle(360.0 - segmentAngleDegrees*2 - segmentAngleDegrees/2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*2 {
			//fmt.Println("multiball hit segment: 2")
			g.balls[i].CalcXYForAngle(360.0 - segmentAngleDegrees - segmentAngleDegrees/2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*3 {
			// fmt.Println("multiball hit segment: 3")
			g.balls[i].CalcXYForAngle(360.0 - segmentAngleDegrees/2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*4 {
			// fmt.Println("multiball hit segment: 4")
			g.balls[i].CalcXYForAngle(segmentAngleDegrees / 2.0)
		} else if ballCenterX < g.paddle.Rect.X+paddleSegmentLenX*5 {
			// fmt.Println("multiball hit segment: 5")
			g.balls[i].CalcXYForAngle(segmentAngleDegrees + segmentAngleDegrees/2.0)
		} else {
			// fmt.Println("multiball hit segment: 6")
			g.balls[i].CalcXYForAngle(segmentAngleDegrees*2 + segmentAngleDegrees/2.0)
		}
		// ensure that the ball bounces upwards
		if g.balls[i].SpeedY > 0 {
			g.balls[i].SpeedY = -g.balls[i].SpeedY
		}

	}
	UpdateNode(g.paddle)
	//g.paddle.Update()

	// update balls
	for i := 0; i < len(g.balls); i++ {
		if g.balls[i] == nil {
			continue
		}
		if g.balls[i].Grabbed {
			g.balls[i].Rect.X = g.paddle.Rect.X
		} else {
			g.balls[i].Update()
		}
	}

	if g.level.TotalHealth <= 0 {
		g.currLevelNum++
		g.level = level.NewLevel(g.currLevelNum)
	}

	fmt.Println(g.balls)
	return nil
}

// inserts a new grabbed ball if max balls not reached
func (g *Game) InitGrabbedBall() {
	// loop to the first available Balls element
	for i := 0; i < config.BallMaxCount; i++ {
		if g.balls[i] == nil {
			g.balls[i] = entities.NewBall(
				g.paddle.Rect.X,
				g.paddle.Rect.Y-config.BallSize,
				config.BallStartingSpeed,
				config.BallStartingAngle,
				true)
			break
		}
	}
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
	game.BallCount = 1

	fmt.Println(game.balls)

	// init paddle
	cursorX, _ := ebiten.CursorPosition()
	game.lives = config.StartingLives
	game.paddle = &entities.Paddle{}
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
