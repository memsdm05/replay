package replay

// Score represents the information seen on the score report
type Score struct {
	// Number of 300s
	N300 uint16

	// Number of 100s in Standard, 150s in Taiko, 100s in CTB, 100s in Mania
	N100 uint16

	// Number of 50s in Standard, small fruit in CTB, 50s in Mania
	N50 uint16

	// Number of Gekis in Standard, Max 300s in Mania
	Gekis uint16

	// Number of Katus in Standard, 200s in Mania
	Katus uint16

	// Number of Misses
	Misses uint16

	// Total score displayed on the score report
	TotalScore int32

	// Greatest combo displayed on the score report
	MaxCombo uint16

	// True if perfect, False if not perfect
	// Being Perfect means the map has been fc'ed and has no dropped slider ends
	IsPerfect bool

	// Mods bitflag
	Mods Mod
}

// TotalHits is the summation of N50, N100, N300, and Misses, cumulating into
// the final amount of note hits.
func (s Score) TotalHits() int {
	return int(s.N50) + int(s.N100) + int(s.N300) + int(s.Misses)
}

func (s Score) Accuracy() float64 {
	return float64(s.N50) * 50 + float64(s.N100) * 100 + float64(s.N300) * 300 /
		float64(s.TotalHits())
}

// Rank returns a Rank enum based off the score
//
// Rank is based off the 2016 client leak. Values may be incorrect.
func (s Score) Rank() Rank {
	r300 := float64(s.N300) / float64(s.TotalHits())
	r50 := float64(s.N50) / float64(s.TotalHits())

	if r300 == 1 {
		if s.Mods & Flashlight & Hidden != 0 {
			return XH
		} else {
			return X
		}
	}

	if r300 > 0.9 && r50 <= 0.01 && s.Misses == 0 {
		if s.Mods & Flashlight & Hidden != 0 {
			return SH
		} else {
			return S
		}
	}

	if (r300 > 0.8 && s.Misses == 0) || (r300 > 0.9) {
		return A
	}
	if (r300 > 0.7 && s.Misses == 0) || (r300 > 0.8) {
		return B
	}
	if r300 > 0.6 {
		return C
	}
	return D
}

func (s *Score) unmarshal(or *osuReader)  {
	var modProxy int32
	or.ReadTypes(
		&s.N300,
		&s.N100,
		&s.N50,
		&s.Gekis,
		&s.Katus,
		&s.Misses,
		&s.TotalScore,
		&s.MaxCombo,
		&s.IsPerfect,
		&modProxy)

	s.Mods = Mod(modProxy)
}

func (s *Score) marshal(ow *osuWriter) {
	ow.WriteTypes(
		s.N300,
		s.N100,
		s.N50,
		s.Gekis,
		s.Katus,
		s.Misses,
		s.TotalScore,
		s.MaxCombo,
		s.IsPerfect,
		int32(s.Mods))
}
