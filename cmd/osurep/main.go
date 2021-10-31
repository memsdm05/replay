package main

import (
	"fmt"
	"github.com/memsdm05/replay"
	"os"
)

func main() {
	f, _ := os.Open("E:\\osu\\Replays\\badeu - Haywyre - Insight [pog] (2020-07-01) Osu.osr")

	r := replay.New(f)

	// smoke takes 8 seconds to fade

	for _, a := range r.Actions {
		if a.KeyState.Has(replay.Smoke) {
			fmt.Printf("%+v\n", a)
		}
	}

	//fmt.Println(r.Actions)
}
