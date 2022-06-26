package match

type MatchStatus int
type Score int

const (
	Idle     MatchStatus = 0
	Playing              = 1
	Finished             = 2
)
