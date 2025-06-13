package priority

type Priority int

const (
	// ColorWhite ColorBrightRed Vivid background-like text highlights for priorities
	ColorWhite        = "\033[97m"
	ColorBrightRed    = "\033[1;91m"
	ColorBrightGreen  = "\033[1;92m"
	ColorBrightYellow = "\033[1;93m"
)

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

func FromInt(code int) Priority {
	switch code {
	case int(Low):
		return Low
	case int(Normal):
		return Normal
	case int(High):
		return High
	default:
		return Priority(-1) // or define an Unknown Priority
	}
}
func ColorForPriority(p string) string {
	switch p {
	case "LOW":
		return ColorBrightYellow // Yellow → attention, but not critical
	case "NORMAL":
		return ColorBrightGreen // Green → okay, default
	case "HIGH":
		return ColorBrightRed // Red → important
	default:
		return ColorWhite
	}
}
