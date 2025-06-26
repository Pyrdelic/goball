package level

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/entities"
)

type Level struct {
	Bricks      []entities.Brick
	TotalHealth int
}

func (l *Level) LoadFromFile(path string) {
	l.Bricks = nil
	// open the level file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		if !(lineCount < config.BrickRowCount) {
			break
		}
		line := scanner.Text()
		fmt.Println(line, len(line), utf8.RuneCountInString(line))
		for columnCount, runeCharacter := range []rune(line) {
			if !(columnCount < config.BrickColumnCount) {
				break // max column count reached
			}
			switch runeCharacter {
			// normal brick
			case 'ä':
				fmt.Println("äääää")
				break
			default:
				// brick undefined
				break
			}
		}
		lineCount++
	}
}
