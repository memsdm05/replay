package replay

type Gamemode int

const (
	Standard = iota
	Taiko
	CTB
	Mania
)

type Button int

const (
	M1 = Button(1 << iota)
	M2
	K1
	K2
	Smoke

	Kb1 = M1 + K1
	kb2 = M2 + K2
)

func (b Button) Has(button Button) bool {
	return b & button != 0
}

type Mod uint32

const(
	None = Mod(0)
	NoFail = Mod(1 << iota)
	Easy
	TouchDevice
	Hidden
	HardRock
	SuddenDeath
	DoubleTime
	Relax
	HalfTime
	Nightcore
	Flashlight
	Autoplay
	SpunOut
	Autopilot
	Perfect
	Key4
	Key5
	Key6
	Key7
	Key8
	FadeIn
	Random
	LastMod
	TargetPractice
	Key9
	Coop
	Key1
	Key3
	Key2
	ScoreV2
	Mirror

	KeyMod = Key4 | Key5 | Key6 | Key7 | Key8
)

type Rank int

const (
	XH = Rank(iota)
	SH
	X
	S
	A
	B
	C
	D
	F
	N
)