package menu

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pyrdelic/goball/button"
	"github.com/pyrdelic/goball/node"
)

type PauseMenu struct {
	ExitGameButton   *button.Button
	ResumeGameButton *button.Button
}

type MainMenu struct {
	ExitGameButton  *button.Button
	StartGameButton *button.Button
}

const (
	ExitGameButtonPressed int = iota + 1
	ResumeButtonPressed
	StartGameButtonPressed
)

func NewPauseMenu() *PauseMenu {
	pm := PauseMenu{}
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
	node.Draw(mm.ExitGameButton, screen)
	node.Draw(mm.StartGameButton, screen)
}

func NewMainMenu() *MainMenu {
	mm := MainMenu{}
	mm.ExitGameButton = button.NewButton(
		100, 100, 50, 50, "Exit",
	)
	mm.StartGameButton = button.NewButton(
		150, 150, 50, 50, "New Game",
	)
	return &mm
}
