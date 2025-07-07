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
	CurrScene    *node.Node
	currLevelNum int
	GameOver     bool
}

func (g *Game) Update() error {

	// update active scene

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// We exit the game by returning a custom error
		return ErrTerminated
	}

	// fmt.Println(g.balls)
	switch node.Update(g.level) {
	// TODO: move to level package
	case level.GameOver:
		fmt.Println("GAME OVER!")
		return ErrTerminated // quit game
	case level.Pause:
		// Switch current scene to pause menu
	default:
		break
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	node.Draw(g.level, screen)

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
