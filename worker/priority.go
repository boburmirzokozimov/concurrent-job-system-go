package worker

type Priority int

const (
	Low Priority = iota
	Normal
	High
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "LOW"
	case Normal:
		return "NORMAL"
	case High:
		return "HIGH"
	default:
		return "UNKNOWN"
	}
}

func ColorForPriority(p Priority) string {
	switch p {
	case Low:
		return ColorBrightYellow // Yellow → attention, but not critical
	case Normal:
		return ColorBrightGreen // Green → okay, default
	case High:
		return ColorBrightRed // Red → important
	default:
		return ColorWhite
	}
}
