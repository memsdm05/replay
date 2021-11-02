package main

import (
	"github.com/memsdm05/replay"
	"os"
)

func main() {
	f, _ := os.Open("E:\\osu\\Replays\\badeu - Haywyre - Insight [pog] (2020-07-01) Osu.osr")
	out, _ := os.Create("test_replay.osr")

	r := replay.New(f)
	//r.Name = "pogdeu"
	r.Marshal(out)

	//fmt.Println(r.Actions)
}
