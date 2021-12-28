package server_engine

var (
	GetSelector Sel = &sel{}
)

type Sel interface {
	Find(sel string) string
}

type sel struct{}

func (s *sel) Find(sel string) string {
	switch sel {
	case "":
		return s.kickass()
	default:
		return ""
	}
}

func (s *sel) kickass() string {
	return ""
}

func (s *sel) gogoAnime() string {
	return ""
}

func (s *sel) nineAnime() string {
	return ""
}

func (s *sel) twistMoe() string {
	return ""
}

func (s *sel) animeheaven() string {
	return ""
}
