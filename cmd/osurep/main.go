package main

import (
	"github.com/memsdm05/replay"
	"math"
	"math/rand"
	"os"
	"time"
)

func main() {
	f, _ := os.Open("E:\\osu\\Replays\\badeu - Haywyre - Insight [pog] (2020-07-01) Osu.osr")
	out, _ := os.Create("test_replay.osr")

	r := replay.New(f)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len(r.Actions); i++ {
		a := &r.Actions[i]
		//a.KeyState |= replay.Smoke
		if !a.KeyState.Has(replay.K1 | replay.K2) {
			a.X	= float64(rand.Intn(584))
			a.Y = float64(rand.Intn(380))

			a.X += math.Cos(float64(i)) * 100
			a.Y += math.Sin(float64(i)) * 100
		}
	}

	r.Name = "pogdeu"
	r.Score.TotalScore = 72727
	r.Timestamp = time.Now()

	r.Marshal(out)

	//fmt.Println(r.Actions)
}
