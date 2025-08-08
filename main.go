package main

import (
	"errors"
	"fmt"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/entities"
	"github.com/pyrdelic/goball/level"
	"github.com/pyrdelic/goball/menu"
	"github.com/pyrdelic/goball/node"
	"github.com/pyrdelic/goball/player"
)

// TODO: Proper event system

type Game struct {
	// TODO: Paddle to level
	// TODO: Balls to level
	//paddle       *entities.Paddle
	balls     [config.BallMaxCount]*entities.Ball
	BallCount int
	//lives        int
	level             *level.Level
	MainMenu          *menu.MainMenu
	PauseMenu         *menu.PauseMenu
	GameOverMenu      *menu.GameOverMenu
	HiScoreMenu       *menu.HiScoreMenu
	CurrScene         node.Node // interface as a type (no * needed)
	currLevelNum      int
	GameOver          bool
	HighScoreAchieved bool

	Player *player.Player
}

//var faceSource *text.GoTextFaceSource

func (g *Game) Update() error {
	if g == nil {
		return nil
	}

	// update active scene
	message := node.Update(g.CurrScene)
	if message.Msg != 0 {
		fmt.Println("TypeStr:", message.TypeStr, "Msg:", message.Msg)
	}
	switch message.TypeStr {
	case "MainMenu":
		switch message.Msg {
		case menu.ExitGameButtonPressed:
			return ErrTerminated // exit game
		case menu.StartGameButtonPressed:
			ebiten.SetCursorMode(ebiten.CursorModeHidden)
			// reset score, lives
			g.Player.Score = 0
			g.level.Lives = config.StartingLives
			// init a new level
			g.level = level.NewLevel(1)
			g.CurrScene = g.level
		}
	case "Level":
		g.Player.Score += message.IntExtra
		switch message.Msg {
		case level.GameOver:
			fmt.Println("GAME OVER")

			g.HiScoreMenu = menu.NewHiScoreMenu(uint64(g.Player.Score))

			ebiten.SetCursorMode(ebiten.CursorModeVisible)
			//g.CurrScene = g.GameOverMenu
			g.CurrScene = g.HiScoreMenu
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
	case "GameOverMenu":
		switch message.Msg {
		case menu.ExitGameButtonPressed:
			fmt.Println("EXITING GAME")
			return ErrTerminated
		case menu.StartGameButtonPressed:
			ebiten.SetCursorMode(ebiten.CursorModeHidden)
			// reset score, lives
			g.Player.Score = 0
			g.level.Lives = config.StartingLives
			// init a new level
			g.level = level.NewLevel(1)
			g.CurrScene = g.level
		}
	case "HiScoreMenu":
		switch message.Msg {
		case menu.MainMenuButtonPressed:
			g.CurrScene = g.MainMenu
		}

	default:
		fmt.Println("Unknown scene")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	node.Draw(g.CurrScene, screen)

	 // // text test
	// str := "asdfasdf"
	// text.Draw(screen, str, &text.GoTextFace{
	// 	Source: faceSource,
	// 	Size:   24,
	// }, &text.DrawOptions{})

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.PlayAreaWidth, config.PlayAreaHeight
}

// Custom error to exit the game loop in a regular way.
var ErrTerminated = errors.New("terminated")

// func initFont() {
// 	face, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	faceSource = face
// }

func main() {

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GO-BALL")
	//initFont()

	game := Game{}
	game.Player = &player.Player{}
	//game.Player.Lives = 3

	// init menus
	game.PauseMenu = menu.NewPauseMenu()
	game.MainMenu = menu.NewMainMenu()
	game.GameOverMenu = menu.NewGameOverMenu()
	game.HiScoreMenu = menu.NewHiScoreMenu(0)

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

	//game.CurrScene = game.MainMenu
	game.CurrScene = game.MainMenu

	ebiten.SetVsyncEnabled(false)
	//ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(&game); err != nil {
		if err == ErrTerminated {
			// Regular termination
			return
		}
		log.Fatal(err)
	}
}
