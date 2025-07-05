package main

import (
	"errors"
	"fmt"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/entities"
	"github.com/pyrdelic/goball/level"
	"github.com/pyrdelic/goball/node"
)

type Game struct {
	// TODO: Paddle to level
	// TODO: Balls to level
	//paddle       *entities.Paddle
	balls     [config.BallMaxCount]*entities.Ball
	BallCount int
	//lives        int
	level        *level.Level
	currLevelNum int
	GameOver     bool
}

// func UpdateNode(n entities.Node) {
// 	n.Update()
// }

// func DrawNode(n entities.Node) {

// }

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
// func speedXYForAngle(speedBase float64, angle float64) (float64, float64) {
// 	radian := angle * (math.Pi / 180)
// 	speedX := speedBase * math.Sin(radian)
// 	speedY := speedBase * math.Cos(radian)
// 	// flip Y component to correct for game space coordinate system
// 	//fmt.Println(int(speedX), int(-speedY))
// 	return speedX, -speedY
// }

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
	//node.UpdateNode(g.level)

	// if g.BallCount <= 0 {
	// 	// lose a life
	// 	g.lives--
	// 	if g.lives >= 0 {
	// 		// initialize a new ball
	// 		g.InitGrabbedBall()
	// 		g.BallCount++
	// 	} else {
	// 		fmt.Println("Out of lives.")
	// 		fmt.Println("GAME OVER")
	// 	}
	// }

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// We exit the game by returning a custom error
		return ErrTerminated
	}

	// fmt.Println(g.balls)
	switch node.Update(g.level) {
	case level.GameOver:
		fmt.Println("GAME OVER!")
		return ErrTerminated // quit game
	default:
		break
	}
	return nil
}

// inserts a new grabbed ball if max balls not reached
// func (g *Game) InitGrabbedBall() {
// 	// loop to the first available Balls element
// 	for i := 0; i < config.BallMaxCount; i++ {
// 		if g.balls[i] == nil {
// 			g.balls[i] = entities.NewBall(
// 				g.paddle.Rect.X,
// 				g.paddle.Rect.Y-config.BallSize,
// 				config.BallStartingSpeed,
// 				config.BallStartingAngle,
// 				true)
// 			break
// 		}
// 	}
// }

func (g *Game) Draw(screen *ebiten.Image) {
	// for iRow := 0; iRow < config.BrickRowCount; iRow++ {
	// 	for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
	// 		g.level.Bricks[iRow][iColumn].Draw(screen)
	// 	}
	// }

	node.Draw(g.level, screen)

	//g.paddle.Draw(screen)

	// for _, b := range g.balls {
	// 	b.Draw(screen)
	// }
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"lives: %d\nlvl health: %d",
		g.level.Lives,
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
	// cursorX, _ := ebiten.CursorPosition()
	//game.lives = config.StartingLives
	// game.paddle = &entities.Paddle{}
	// game.paddle.Rect.X = float64(cursorX)
	// game.paddle.Rect.Y = 200
	// game.paddle.Rect.W = config.PaddleStartingWidth
	// game.paddle.Rect.H = 5
	// game.paddle.Image = ebiten.NewImage(int(game.paddle.Rect.W), int(game.paddle.Rect.H))
	// game.paddle.Image.Fill(color.White)

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
