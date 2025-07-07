package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pyrdelic/goball/button"
	"github.com/pyrdelic/goball/node"
)

type PauseMenu struct {
	ExitGameButton   *button.Button
	ResumeGameButton *button.Button
}

const (
	ExitGameButtonPressed node.Message = iota + 1
	ResumeButtonPressed
)

func NewPauseMenu() *PauseMenu {
	pm := PauseMenu{}
	pm.ExitGameButton = button.NewButton(100, 100, 30, 30)
	pm.ResumeGameButton = button.NewButton(150, 150, 30, 30)
	return &pm
}

func (pm *PauseMenu) Update() node.Message {
	if pm == nil {
		return 0
	}
	if pm.ResumeGameButton.IsJustClicked() {
		return ResumeButtonPressed
	}
	if pm.ExitGameButton.IsJustClicked() {
		return ExitGameButtonPressed
	}
	return 0
}

func (pm *PauseMenu) Draw(screen *ebiten.Image) {
	if pm == nil {
		return
	}
	pm.ExitGameButton.Draw(screen)
	pm.ResumeGameButton.Draw(screen)
}
