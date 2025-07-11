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
	"github.com/pyrdelic/goball/menu"
	"github.com/pyrdelic/goball/node"
	"github.com/pyrdelic/goball/player"
)

type Game struct {
	// TODO: Paddle to level
	// TODO: Balls to level
	//paddle       *entities.Paddle
	balls     [config.BallMaxCount]*entities.Ball
	BallCount int
	//lives        int
	level        *level.Level
	PauseMenu    *menu.PauseMenu
	CurrScene    node.Node
	currLevelNum int
	GameOver     bool

	Player *player.Player
}

func (g *Game) Update() error {

	// update active scene

	message := node.Update(g.CurrScene)
	if message.Msg != 0 {
		fmt.Println("TypeStr:", message.TypeStr, "Msg:", message.Msg)
	}
	switch message.TypeStr {
	case "Level":
		g.Player.Score += message.IntExtra
		switch message.Msg {
		case level.GameOver:
			fmt.Println("GAME OVER")
			// TODO: Game over / hi-score scene
			return ErrTerminated // exit game
		case level.Pause:
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
			fmt.Println("PAUSE")
			g.CurrScene = g.PauseMenu
			// TODO: switch current scene to PauseMenu
		}
	case "PauseMenu":
		switch message.Msg {
		case menu.ExitGameButtonPressed:
			fmt.Println("EXITING GAME")
			return ErrTerminated
		case menu.ResumeButtonPressed:
			// TODO: Switch current scene back to level
			fmt.Println("RESUMING")
			ebiten.SetCursorMode(ebiten.CursorModeHidden)
			g.CurrScene = g.level
		}
	default:
		fmt.Println("Unknown scene")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	node.Draw(g.CurrScene, screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"lives: %d\nscore: %d\nlvl health: %d",
		g.level.Lives,
		g.Player.Score,
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
	game.Player = &player.Player{}
	//game.Player.Lives = 3

	// init pause menu
	game.PauseMenu = menu.NewPauseMenu()

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

	game.CurrScene = game.level

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
