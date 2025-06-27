package level

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/entities"
)

type Level struct {
	Bricks      [config.BrickRowCount][config.BrickColumnCount]*entities.Brick
	TotalHealth int
}

func (l *Level) LoadFromFile(path string) {
	//l.Bricks = nil
	// open the level file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	iRow := 0
	for scanner.Scan() {
		if !(iRow < config.BrickRowCount) {
			break
		}
		line := scanner.Text()
		fmt.Println(line, len(line), utf8.RuneCountInString(line))
		for iColumn, runeCharacter := range []rune(line) {
			if !(iColumn < config.BrickColumnCount) {
				break // max column count reached
			}
			switch runeCharacter {
			case '0':
				// no brick
				l.Bricks[iRow][iColumn] = nil
				fmt.Println("No brick")
			case '1':
				// basic brick

				fmt.Println("Basic brick")
				brick := entities.Brick{}
				brick.Health = 1
				brick.BrickType = 1
				brick.Rect.X = float64(iColumn * config.BrickWidth)
				brick.Rect.Y = float64(iRow * config.BrickHeight)
				brick.Rect.W = config.BrickWidth
				brick.Rect.H = config.BrickHeight
				brick.Image = ebiten.NewImage(
					int(brick.Rect.W),
					int(brick.Rect.H))
				brick.Image.Fill(color.RGBA{
					R: uint8(64),
					G: uint8(255),
					B: uint8(64),
					A: uint8(255)})
				l.TotalHealth += brick.Health
				l.Bricks[iRow][iColumn] = &brick

				// default to no brick
			default:
				// default to no brick
				fmt.Println("No brick")
				l.Bricks[iRow][iColumn] = nil
			}
		}
		iRow++
	}
}

// returns a pointer to a new level, based on the level number.
func NewLevel(levelNumber int) *Level {
	level := Level{}
	levelPath := fmt.Sprintf("levels/level%d.txt", levelNumber)
	fmt.Println("Loading level from file:", levelPath)
	level.LoadFromFile(levelPath)
	return &level
}

func (l *Level) PrintLevel() {
	for iRow := 0; iRow < config.BrickRowCount; iRow++ {
		for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
			var brickType int
			if l.Bricks[iRow][iColumn] == nil {
				brickType = 0
			} else {
				brickType = l.Bricks[iRow][iColumn].BrickType
			}
			fmt.Printf("%d", brickType)
		}
		fmt.Println()
	}
}
