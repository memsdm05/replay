package main

import (
	"github.com/memsdm05/replay"
	"os"
	"time"
)

func main() {
	f, _ := os.Open("E:\\osu\\Replays\\badeu - Haywyre - Insight [pog] (2020-07-01) Osu.osr")
	out, _ := os.Create("test_replay.osr")

	r := replay.New(f)


	r.Name = "\n\n\n\n\n\n\n                                       loled"
	r.Timestamp = time.Unix(459829320, 0)
	r.Score.MaxCombo = 727
	r.Score.TotalScore = 7272727
	r.Marshal(out)

	//fmt.Println(r.Actions)
}
