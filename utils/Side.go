package utils

// Strategy Side.
type Side uint32

type Sides []Side

type SideData struct {
	Name string
	Data []Side
}

const (
	Buy  Side = 0
	Sell Side = 1 << iota
	Hold
)

// String implements the stringer interface
func (s Side) String() string {
	switch s {
	case Buy:
		return "BUY"
	case Sell:
		return "SELL"
	case Hold:
		return "HOLD"
	default:
		return "UNKNOWN"
	}
}

// String implements the stringer interface
func (s Sides) Strings() []string {
	strs := make([]string, len(s))
	for x := range s {
		strs[x] = s[x].String()
	}
	return strs
}
