package parse

type ActionType int

const (
	Shift ActionType = iota
	Accept
	Reduce
	Error
	Conflict       // Was a reduce, but part of a conflict
	ShiftResolved  // Was a shift. Precedence resolved conflict
	ReduceResolved // Was reduce. Precedence resolved conflict
	NotUsed        // Deleted by compression
)

// Every shift or reduce operation is stored as one of the following.
type Action struct {
}
