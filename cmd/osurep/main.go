package main

import (
	"github.com/memsdm05/replay"
	"math/rand"
	"os"
	"time"
)

func main() {
	f, _ := os.Open("E:\\osu\\Replays\\badeu - Haywyre - Insight [pog] (2020-07-01) Osu.osr")
	out, _ := os.Create("test_replay.osr")

	r := replay.New(f)
	rand.Seed(time.Now().UnixNano())

	min, max := float64(-10), float64(10)

	for i := 0; i < len(r.Actions); i++ {
		a := &r.Actions[i]
		a.KeyState |= replay.Smoke
		a.X += min + rand.Float64() * (max - min)
		a.Y += min + rand.Float64() * (max - min)
	}

	r.Name = "pogdeu"
	//r.Timestamp = time.Now()

	r.Marshal(out)

	//fmt.Println(r.Actions)
}
