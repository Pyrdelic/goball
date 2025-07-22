package menu

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pyrdelic/goball/button"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/hiscore"
	"github.com/pyrdelic/goball/node"
)

type HiScoreMenu struct {
	Title          string
	HiScores       *[config.HiScoreTopCount]hiscore.HiScore
	HiScoresStr    string
	MainMenuButton *button.Button
}

func (hsm *HiScoreMenu) Update() node.Message {

	if hsm == nil || hsm.HiScores == nil {
		return node.Message{TypeStr: "HiScoreMenu"}
	}
	//fmt.Println("Hiscoremenu update")
	if hsm.MainMenuButton.IsJustClicked() {
		return node.Message{
			TypeStr: "HiScoreMenu",
			Msg:     MainMenuButtonPressed,
		}
	}
	return node.Message{TypeStr: "HiScoreMenu"}
}

func (hsm *HiScoreMenu) Draw(screen *ebiten.Image) {
	if hsm == nil {
		return
	}
	//fmt.Println("HiScoremenu draw")
	drawTitleText(hsm.Title, screen)

	// TODO: draw hi scores
	if hsm.HiScores != nil {
		textop := text.DrawOptions{}
		textop.GeoM.Translate(float64(35), float64(50))
		textop.LineSpacing = 14
		textop.ColorScale.ScaleWithColor(color.White)
		text.Draw(
			screen,
			hsm.HiScoresStr,
			&text.GoTextFace{
				Source: contentTextFaceSource,
				Size:   14,
			},
			&textop)
	}

	node.Draw(hsm.MainMenuButton, screen)
}

func NewHiScoreMenu() *HiScoreMenu {
	initContentTextFaceSource()
	hsm := HiScoreMenu{}
	hsm.HiScores = &[config.HiScoreTopCount]hiscore.HiScore{}
	hsm.Title = "Hi-Scores"
	hiscore.LoadHiScores(hsm.HiScores, "hiscore/hiscore.txt")

	hsm.HiScoresStr = ""
	for i := range len(hsm.HiScores) {
		if i > config.HiScoreTopCount {
			break
		}
		row := strings.Join(
			[]string{hsm.HiScores[i].Name, fmt.Sprintf("%d", hsm.HiScores[i].Score)},
			" ",
		)
		hsm.HiScoresStr += row + "\n"
	}
	fmt.Println(hsm.HiScoresStr)

	hsm.MainMenuButton = button.NewButton(
		config.PlayAreaWidth/4*3,
		config.PlayAreaHeight/6*5,
		50,
		20,
		"Main Menu",
	)
	return &hsm
}
