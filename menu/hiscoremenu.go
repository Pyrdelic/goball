package menu

import (
	"fmt"
	"image/color"
	"strings"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pyrdelic/goball/button"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/hiscore"
	"github.com/pyrdelic/goball/node"
)

type HiScoreMenu struct {
	Title                  string
	HiScores               *[config.HiScoreTopCount]hiscore.HiScore
	HiScoresStr            string
	MainMenuButton         *button.Button
	HiScoreToSubmit        hiscore.HiScore
	HiScoreToSubmitIsNewHS bool
	nameInputEnabled       bool
	newHSPosition          int
}

var (
	path string = "hiscore/hiscore.txt"
)

// returns -1 if out of scoreboard
func (hsm *HiScoreMenu) getScoreBoardPosition(score uint64) int {
	// get position on score board
	//fmt.Println(score)
	scoreBoardPosition := -1
	for i := range len(hsm.HiScores) {
		//fmt.Println(hsm.HiScoreToSubmit.Score)
		if scoreBoardPosition == -1 && score > hsm.HiScores[i].Score {
			//fmt.Println("in if")
			scoreBoardPosition = i
			hsm.HiScoreToSubmitIsNewHS = true
			//fmt.Println("ScoreBoardPosition:", scoreBoardPosition)
			hsm.nameInputEnabled = true
			hsm.newHSPosition = scoreBoardPosition
			break
		}
	}
	return scoreBoardPosition

}

// Inserts a new HiScore into HiScores,
// dropping those below config.HiScoreTopCount
func (hsm *HiScoreMenu) insertAndTrucateHighScore(slice []hiscore.HiScore, score uint64, pos int) *[]hiscore.HiScore {
	fmt.Println(pos)
	if pos < 0 || pos > config.HiScoreTopCount {
		panic("HiScore index out of bounds")
	}

	slice = append(slice, hiscore.HiScore{})

	copy(slice[pos+1:], slice[pos:])

	slice[pos] = hiscore.HiScore{Name: "", Score: score}
	if len(slice) > config.HiScoreTopCount {
		slice = slice[:config.HiScoreTopCount]
	} else {
		slice = slice[:config.HiScoreTopCount-1]
	}
	return &slice
}

func (hsm *HiScoreMenu) Update() node.Message {
	message := node.Message{TypeStr: "HiScoreMenu"}
	if hsm == nil || hsm.HiScores == nil {
		return message
	}
	// Button presses
	if hsm.MainMenuButton.IsJustClicked() {
		message.Msg = MainMenuButtonPressed
		hiscore.WriteHiScores(hsm.HiScores, path)
		return message
	}

	// Name input
	if hsm.nameInputEnabled {
		readChars := []rune{}
		readChars = ebiten.AppendInputChars(readChars)
		for i := range len(readChars) {
			if !(utf8.RuneCount([]byte(hsm.HiScoreToSubmit.Name)) < 3) {
				break
			}
			hsm.HiScores[hsm.newHSPosition].Name += string(readChars[i])
		}
		hsm.HiScoresStr = hiScoresToStr(hsm.HiScores)
	}

	return message

	// // A very shoddy way of getting possible
	// // hiscore name input. Only returns the first
	// // key pressed during a tick, discards rest.
	// message := node.Message{TypeStr: "HiScoreMenu"}
	// inputtedChars := []rune{}
	// inputtedChars = ebiten.AppendInputChars(inputtedChars)
	// if len(inputtedChars) > 0 {
	// 	message.Msg = KeyPressed
	// 	message.IntExtra = int(inputtedChars[0])
	// }

	// return message
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

func hiScoresToStr(hiScores *[config.HiScoreTopCount]hiscore.HiScore) string {
	hiScoresStr := ""
	for i := range len(hiScores) {
		if i > config.HiScoreTopCount {
			break
		}
		row := strings.Join(
			[]string{hiScores[i].Name, fmt.Sprintf("%d", hiScores[i].Score)},
			" ",
		)
		hiScoresStr += row + "\n"
	}
	return hiScoresStr
}

// 0 as scoreToSubmit considered as no submission.
func NewHiScoreMenu(scoreToSubmit uint64) *HiScoreMenu {
	fmt.Println("NewHighScoreMenu: ", scoreToSubmit)
	initContentTextFaceSource()
	hsm := HiScoreMenu{}
	// hsm.HiScoreToSubmitIsNewHS = false
	// hsm.HiScoreToSubmit = hiscore.HiScore{Name: "", Score: scoreToSubmit}
	hsm.HiScores = &[config.HiScoreTopCount]hiscore.HiScore{}
	hsm.Title = "Hi-Scores"
	hiscore.LoadHiScores(hsm.HiScores, path)
	if scoreToSubmit != 0 {
		pos := hsm.getScoreBoardPosition(scoreToSubmit)
		fmt.Println(pos)
		if pos != -1 {
			hsm.HiScores = (*[10]hiscore.HiScore)(*hsm.insertAndTrucateHighScore(hsm.HiScores[:], scoreToSubmit, pos))
			hsm.HiScoreToSubmitIsNewHS = true
		}

	}

	hsm.HiScoresStr = hiScoresToStr(hsm.HiScores)

	//fmt.Println(hsm.HiScoresStr)

	hsm.MainMenuButton = button.NewButton(
		config.PlayAreaWidth/4*3,
		config.PlayAreaHeight/6*5,
		50,
		20,
		"Main Menu",
	)
	return &hsm
}
