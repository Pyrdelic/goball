package hiscore_test

import (
	"testing"

	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/hiscore"
)

const (
	path = "hiscore.txt"
)

// func TestLoadHiScore(t *testing.T) {
// 	hiScores := [config.HiScoreTopCount]hiscore.HiScore{}
// 	hiscore.LoadHiScores(&hiScores, path)
// 	if hiScores[0].Name != "PYY" {
// 		t.Errorf("hiScores[0].Name == \"%s\", expected PYY", hiScores[0].Name)
// 	}
// 	if hiScores[0].Score != 900001 {
// 		t.Errorf("hiScores[0].Score == %d, expected 900001", hiScores[0].Score)
// 	}
// }

func TestWriteHiScore(t *testing.T) {
	path := "hiscore.txt"
	hiScores := [config.HiScoreTopCount]hiscore.HiScore{
		{Name: "ASD", Score: 9001},
		{Name: "LED", Score: 700},
		{Name: "LCD", Score: 5000},
		{Name: "WAV", Score: 3500},
	}
	//hiscore.LoadHiScores(&hiScores, path)
	//hiScores[3] = hiscore.HiScore{Name: "JEE", Score: 1337}
	hiscore.WriteHiScores(&hiScores, path)
	hiscore.LoadHiScores(&hiScores, path)
	if hiScores[3].Name != "LED" {
		t.Errorf("hiScores[3].Name == \"%s\", expected LED", hiScores[3].Name)
	}
}
