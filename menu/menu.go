package menu

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pyrdelic/goball/button"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/node"
)

const (
	ExitGameButtonPressed int = iota + 1
	ResumeButtonPressed
	StartGameButtonPressed
	MainMenuButtonPressed
	KeyPressed
)

type PauseMenu struct {
	Text             string
	ExitGameButton   *button.Button
	ResumeGameButton *button.Button
}

type MainMenu struct {
	Text            string
	ExitGameButton  *button.Button
	StartGameButton *button.Button
}

type Menu struct {
	//typeStr string
	Title   string
	Buttons []*button.Button
}

// Common text face for menu title
var titleTextFaceSource *text.GoTextFaceSource
var titleTextFaceSourceInitialized bool = false
var contentTextFaceSource *text.GoTextFaceSource
var contentTextFaceSourceInitialized bool = false

// init common text for buttons
func initTitleTextFaceSource() {
	if !titleTextFaceSourceInitialized {
		face, err := text.NewGoTextFaceSource(
			bytes.NewReader(fonts.MPlus1pRegular_ttf))
		if err != nil {
			log.Fatal(err)
		}
		titleTextFaceSource = face
		titleTextFaceSourceInitialized = true
	}
}

func initContentTextFaceSource() {
	if !contentTextFaceSourceInitialized {
		face, err := text.NewGoTextFaceSource(
			bytes.NewReader(fonts.MPlus1pRegular_ttf))
		if err != nil {
			log.Fatal(err)
		}
		contentTextFaceSource = face
		titleTextFaceSourceInitialized = true
	}
}

func drawTitleText(titleText string, screen *ebiten.Image) {
	textop := text.DrawOptions{}
	textop.GeoM.Translate(float64(config.MenuTitleTextX), float64(config.MenuTitleTextY))
	textop.ColorScale.ScaleWithColor(color.White)
	text.Draw(
		screen,
		titleText,
		&text.GoTextFace{
			Source: titleTextFaceSource,
			Size:   20,
		},
		&textop)
}

// // NewMenu returns a pointer to a new empty Menu.
// func NewMenu(typeStr string) *Menu {
// 	m := Menu{}
// 	m.TypeStr = typeStr
// 	return &m
// }

// func (m *Menu) Update() node.Message {
// 	if m == nil {
// 		return node.Message{}
// 	}

// 	for i := range len(m.Buttons) {
// 		event := node.Update(m.Buttons[i])

// 	}
// 	return node.Message{TypeStr: m.Title}
// }

// func (m *Menu) Draw(screen *ebiten.Image) {

// }

func NewPauseMenu() *PauseMenu {
	initTitleTextFaceSource()
	pm := PauseMenu{}
	pm.Text = "Pause"
	pm.ExitGameButton = button.NewButton(100, 100, 30, 30, "Exit")
	pm.ResumeGameButton = button.NewButton(150, 150, 30, 30, "Resume")
	return &pm
}

// Pause menu
func (pm *PauseMenu) Update() node.Message {
	if pm == nil {
		return node.Message{
			TypeStr: "nil",
			Msg:     0,
		}
	}
	if pm.ResumeGameButton.IsJustClicked() {
		fmt.Println("Resume just clicked")
		return node.Message{
			TypeStr: "PauseMenu",
			Msg:     ResumeButtonPressed,
		}
	}
	if pm.ExitGameButton.IsJustClicked() {
		return node.Message{
			TypeStr: "PauseMenu",
			Msg:     ExitGameButtonPressed,
		}
	}
	return node.Message{
		TypeStr: "PauseMenu",
		Msg:     0,
	}
}

func (pm *PauseMenu) Draw(screen *ebiten.Image) {
	if pm == nil {
		return
	}
	drawTitleText(pm.Text, screen)
	pm.ExitGameButton.Draw(screen)
	pm.ResumeGameButton.Draw(screen)
}

// Main menu
func (mm *MainMenu) Update() node.Message {
	if mm == nil {
		return node.Message{
			TypeStr: "MainMenu",
		}
	}
	if mm.ExitGameButton.IsJustClicked() {
		return node.Message{
			TypeStr: "MainMenu",
			Msg:     ExitGameButtonPressed,
		}
	}
	if mm.StartGameButton.IsJustClicked() {
		return node.Message{
			TypeStr: "MainMenu",
			Msg:     StartGameButtonPressed,
		}
	}
	return node.Message{
		TypeStr: "MainMenu",
	}
}

func (mm *MainMenu) Draw(screen *ebiten.Image) {
	drawTitleText(mm.Text, screen)
	node.Draw(mm.ExitGameButton, screen)
	node.Draw(mm.StartGameButton, screen)
}

func NewMainMenu() *MainMenu {
	initTitleTextFaceSource()
	mm := MainMenu{}
	mm.Text = "GO-BALL"
	mm.ExitGameButton = button.NewButton(
		100, 100, 50, 50, "Exit",
	)
	mm.StartGameButton = button.NewButton(
		150, 150, 50, 50, "New Game",
	)
	return &mm
}

type GameOverMenu struct {
	Text           string
	NewGameButton  *button.Button
	ExitGameButton *button.Button
}

func (gom *GameOverMenu) Update() node.Message {
	if gom == nil {
		return node.Message{TypeStr: "GameOverMenu"}
	}
	if gom.ExitGameButton.IsJustClicked() {
		return node.Message{
			TypeStr: "GameOverMenu",
			Msg:     ExitGameButtonPressed,
		}
	}
	if gom.NewGameButton.IsJustClicked() {
		return node.Message{
			TypeStr: "GameOverMenu",
			Msg:     StartGameButtonPressed,
		}
	}
	return node.Message{TypeStr: "GameOverMenu"}
}

func (gom *GameOverMenu) Draw(screen *ebiten.Image) {
	if gom == nil {
		return
	}
	drawTitleText(gom.Text, screen)
	node.Draw(gom.ExitGameButton, screen)
	node.Draw(gom.NewGameButton, screen)

}

func NewGameOverMenu() *GameOverMenu {
	initTitleTextFaceSource()
	gom := GameOverMenu{}
	gom.Text = "Game over."
	gom.NewGameButton = button.NewButton(100, 100, config.ButtonWidth, config.ButtonHeight, "New Game")
	gom.ExitGameButton = button.NewButton(150, 150, config.ButtonWidth, config.ButtonHeight, "Exit Game")
	return &gom
}
