package server_engine

var (
	GetSelector Sel = &sel{}
)

type Sel interface {
	Find() string
}

type sel struct{}

func (s *sel) Find() string {
	return ""
}

func (s *sel) kickass() string {
	return ""
}
