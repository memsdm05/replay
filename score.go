package replay

type Score struct {
	N300 uint16
	N100 uint16
	N50 uint16
	Gekis uint16
	Katus uint16
	Misses uint16
	TotalScore int32
	MaxCombo uint16
	IsPerfect bool
	Mods Mod
}

func (s Score) TotalHits() uint16 {
	return s.N50 + s.N100 + s.N300 + s.Misses
}

func (s Score) Accuracy() float64 {
	return float64(s.N50) * 50 + float64(s.N100) * 100 + float64(s.N300) * 300 /
		float64(s.TotalHits())
}

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
