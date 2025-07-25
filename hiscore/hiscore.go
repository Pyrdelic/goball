package hiscore

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/pyrdelic/goball/config"
)

type HiScore struct {
	Name  string
	Score uint64
}

// // Named type for a slice of HiScore
// type ByScore []HiScore

// // sort.Interface implementation
// func (hs ByScore) Len() int           { return len(hs) }
// func (hs ByScore) Swap(i, j int)      { hs[i], hs[j] = hs[j], hs[i] }
// func (hs ByScore) Less(i, j int) bool { return hs[i].Score < hs[j].Score }

func cmp(a HiScore, b HiScore) int {
	if a.Score < b.Score {
		return -1
	} else if b.Score < a.Score {
		return 1
	}
	return 0
}

func LoadHiScores(hiScores *[config.HiScoreTopCount]HiScore, path string) {
	if hiScores == nil {
		return
	}

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
	slices.SortFunc(hiScores[:], cmp)
	slices.Reverse(hiScores[:])
	// for i := range len(hiScores) {
	// 	fmt.Println(hiScores[i].Name, hiScores[i].Score)
	// }
}
func WriteHiScores(hiScores *[config.HiScoreTopCount]HiScore, path string) {
	if hiScores == nil {
		return
	}
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()
	slices.SortFunc(hiScores[:], cmp)
	slices.Reverse(hiScores[:])

	for i := range len(hiScores) {
		if !(i < config.HiScoreTopCount) {
			break
		}
		line := fmt.Sprintf("%s %d", hiScores[i].Name, hiScores[i].Score)
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}
