package hiscore

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pyrdelic/goball/config"
)

type HiScore struct {
	Name  string
	Score uint64
}

func LoadHiScores(hiScores *[config.HiScoreTopCount]HiScore, path string) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	hiScoresRead := 0
	for scanner.Scan() {
		if hiScoresRead >= 10 {
			break
		}

		line := scanner.Text()
		elements := strings.Split(line, " ")

		// validate
		if len(elements) < 2 {
			continue
		}

		if score, err := strconv.ParseUint(elements[1], 10, 64); err == nil {
			hiScores[hiScoresRead] = HiScore{Name: elements[0], Score: score}
			hiScoresRead++
		}
		continue
	}
}
func WriteHiScores(hiScores *[config.HiScoreTopCount]HiScore, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for i := range len(hiScores) {
		line := fmt.Sprintf("%s %d", hiScores[i].Name, hiScores[i].Score)
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}
