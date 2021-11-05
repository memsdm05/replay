package replay

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Replay Represents an osu! .osr file. All integers are sized, while other fields are abstracted.
type Replay struct {
	// Mode is the Gamemode that the replay is played in
	Mode Gamemode

	// Version is the game version that the replay was made in
	Version int32

	// BeatmapDigest is the MD5 hash of the replay beatmap
	BeatmapDigest [md5.Size]byte

	// Name is the player who supposedly created this replay
	Name string

	// LoadedHash is what the replay file had as its replay hash
	LoadedHash string

	// Score is the score screen
	Score Score

	// LifeBarGraph is a slice of life, each representing a point on the graph
	LifeBarGraph []Life

	// Timestamp is when this replay was made
	Timestamp time.Time

	// Actions is a slice of cursor positions and states, in the form of an Action
	Actions []Action

	// ScoreId is the associated online score id if applicable
	ScoreId int64

	// When PureTS is false, random nanoseconds are added to the marshalled replay
	// so that osu! can play different replays with the same timestamp.
	//
	// By default, PureTs is false.
	PureTS bool
}

// New creates a new Replay from an io.Reader
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
	reader.ReadTypes(&ret.Name, &ret.LoadedHash)
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

// Marshal writes the Replay to an io.Writer in .osr format.
// As of v0.0.2, this function creates replays that can be read by osu
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
	//writer.WriteTypes("")

	var ts = r.Timestamp
	if !r.PureTS {
		ran := rand.New(rand.NewSource(time.Now().UnixNano()))
		ts.Add(time.Duration(ran.Intn(10000000)))
	}
	writer.WriteTypes(ts.Unix() * 10000000 + 621355968000000000) // UnixNano -> .NET ticks

	//writer.WriteTypes(int32(0))
	r.createActions(writer)
	writer.WriteTypes(r.ScoreId)
}

func (r *Replay) createGraph(ow *osuWriter) {
	var sb strings.Builder
	for _, l := range r.LifeBarGraph {
		sb.WriteString(l.entry() + ",")
	}
	ow.WriteTypes(sb.String())
}

func must(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func (r *Replay) createActions(ow *osuWriter) {
	var buf bytes.Buffer
	stream, _ := lzma.NewWriter(&buf)
	for _, a := range r.Actions {
		stream.Write([]byte(a.entry() + ","))
	}
	stream.Close()

	fmt.Println(buf.Len())
	ow.WriteTypes(int32(buf.Len()))
	buf.WriteTo(ow)
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
			if len(r.Actions) > 0 {
				a.Offset = r.Actions[len(r.Actions) - 1].Offset + a.Since
			} else {
				a.Offset = a.Since
			}
			r.Actions = append(r.Actions, a)
			accum.Reset()
		}
	}
}

// Hash is a MD5 digest of selected fields. Used for cheat detection and replay validation.
func (r *Replay) Hash() [md5.Size]byte {
	var perfectStr string
	if r.Score.IsPerfect {
		perfectStr = "True"
	} else {
		perfectStr = "False"
	}

	s := fmt.Sprintf("%vp%vo%vo%vt%va%xr%ve%vy%vo%vu0%vTrue",
		r.Score.N100 + r.Score.N300,
		r.Score.N50,
		r.Score.Gekis,
		r.Score.Katus,
		r.Score.Misses,
		r.BeatmapDigest,
		r.Score.MaxCombo,
		perfectStr,
		r.Name,
		r.Score.TotalScore,
		r.Score.Mods)

	return md5.Sum([]byte(s))
}

// Life is a single entry on the LifeGraph
type Life struct {
	// Health
	Health float64
	Offset time.Duration
}

// Entry is how Life is represented in a .osr file
func (l Life) entry() string {
	return fmt.Sprintf("%f|%d", l.Health, l.Offset.Milliseconds())
}

// Action represents a single replay frame
type Action struct {
	// X and Y is the position of the cursor
	X, Y float64

	// KeyState is current button press state
	// See Button for more information
	KeyState Button

	// Since is the time.Duration since the last action
	Since time.Duration

	// Offset is a time.Duration of where the action is in
	// relation to the start of the map
	Offset time.Duration
}

// Entry is how Action is represented in the LZMA stream
func (a Action) entry() string {
	return fmt.Sprintf("%d|%f|%f|%d", a.Since.Milliseconds(), a.X, a.Y, a.KeyState)
}
