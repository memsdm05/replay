package replay

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"strconv"
	"strings"
	"time"
)

type Replay struct {
	Mode Gamemode
	Version int32
	BeatmapDigest [md5.Size]byte
	Name string // player name
	hash string
	Score Score
	LifeBarGraph []Life
	Timestamp time.Time // numOfTicks / 10000 - 62136892800000
	Actions []Action
	ScoreId int64
}


func New(r io.Reader) *Replay {
	ret := new(Replay)

	reader := &osuReader{ r }

	// mode and version
	var mode byte
	reader.ReadTypes(&mode, &ret.Version)
	ret.Mode = Gamemode(mode)

	// beatmap digest
	var digestS string
	reader.ReadTypes(&digestS)
	digest, _ := hex.DecodeString(digestS)
	copy(ret.BeatmapDigest[:], digest)

	// name, replay hash, life graph, and score frame
	reader.ReadTypes(&ret.Name, &ret.hash)
	ret.Score.unmarshal(reader)
	ret.parseGraph(reader)

	// timestamp and actions
	var ticks int64
	reader.ReadTypes(&ticks)
	ret.Timestamp = time.Unix((ticks - 621355968000000000) / 10000000, 0) // .NET ticks -> UNIX
	ret.parseActions(reader)

	// score id
	reader.ReadTypes(&ret.ScoreId)

	return ret
}

func (r *Replay) Marshal(w io.Writer)  {
	writer := &osuWriter{ w }

	writer.WriteTypes(
		byte(r.Mode),
		r.Version,
		fmt.Sprintf("%x", r.BeatmapDigest),
		r.Name,
		fmt.Sprintf("%x", r.Hash()))

	r.Score.marshal(writer)
	r.createGraph(writer)

	// todo fix conversion
	writer.WriteTypes(int64(r.Timestamp.Unix() + 621355968000000000) / 100) // UnixNano -> .NET ticks

	r.createActions(writer)
	writer.WriteTypes(r.ScoreId)
}

func (r *Replay) createGraph(ow *osuWriter) {
	var sb strings.Builder
	for _, l := range r.LifeBarGraph {
		sb.WriteString(l.String())
		sb.WriteRune(',')
	}
}

func (r *Replay) createActions(ow *osuWriter) {
	var buf bytes.Buffer
	for _, a := range r.Actions {
		buf.WriteString(a.String())
		buf.WriteRune(',')
	}
	ow.WriteTypes(int32(buf.Len()))
	stream, _ := lzma.NewWriter(ow)
	io.Copy(stream, &buf)
}

func (r *Replay) parseGraph(or *osuReader)  {
	var graphS string
	or.ReadTypes(&graphS)
	for _, pair := range strings.Split(graphS, ",") {
		if pair == "" { continue }

		var l Life
		split := strings.Split(pair, "|")

		l.Health, _ = strconv.ParseFloat(split[1], 32)
		offsetInt, _ := strconv.Atoi(split[0])
		l.Offset = time.Duration(offsetInt) * time.Millisecond

		r.LifeBarGraph = append(r.LifeBarGraph, l)
	}
}

func (r *Replay) parseActions(or *osuReader)  {
	var actionSize uint32
	or.ReadTypes(&actionSize)
	t, _ := lzma.NewReader(io.LimitReader(or, int64(actionSize))) // what the hell is this shit
	stream := bufio.NewReader(t)

	var accum bytes.Buffer
	for {
		b, e := stream.ReadByte()
		if e == io.EOF {
			break
		}
		accum.WriteByte(b)

		if b == byte(',') {
			var a Action
			var tmp struct{
				since int
				state int
			}

			values := bytes.Split(accum.Bytes()[:accum.Len() - 1], []byte{'|'})
			if len(values) != 4 { continue }

			tmp.since, _ = strconv.Atoi(string(values[0]))
			a.X, _ = strconv.ParseFloat(string(values[1]), 32)
			a.Y, _ = strconv.ParseFloat(string(values[2]), 32)
			tmp.state, _ = strconv.Atoi(string(values[3]))

			a.Since = time.Duration(tmp.since) * time.Millisecond
			a.KeyState = Button(tmp.state)

			r.Actions = append(r.Actions, a)
			accum.Reset()
		}
	}
}

func (r *Replay) Hash() [md5.Size]byte {
	s := fmt.Sprintf("%dosu%s%x%d%d",
		r.Score.MaxCombo,
		r.Name,
		r.BeatmapDigest,
		r.Score.TotalScore,
		r.Score.Rank())

	return md5.Sum([]byte(s))
}

type Life struct {
	Health float64
	Offset time.Duration
}

func (l Life) String() string {
	return fmt.Sprintf("%f|%d", l.Health, l.Offset / time.Millisecond)
}

type Action struct {
	X, Y float64
	KeyState Button
	Since time.Duration
}

func (a Action) String() string {
	return fmt.Sprintf("%d|%f|%f|%d", a.Since.Milliseconds(), a.X, a.Y, a.KeyState)
}


