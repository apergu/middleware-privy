package constants

const (
	StatusInActive int8 = 1 << iota
	StatusActive
	StatusPublished
)

func FindStatus(src, dst int8) bool {
	return src&dst == dst
}
