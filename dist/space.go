package dist

// An ID for a particular outcome in a space
type Outcome int

// A set of outcomes for some probability measure.
type Space interface{}

// A subset of the reals
type RealSpace interface {
	Space

	// The infimum (min) value in the space, or negative infinity
	Inf() float64

	// The supremum (max) value in the space, or positive infinity
	Sup() float64
}

// A sample space over a discrete set
type DiscreteSpace interface {

	// Every discrete space is a space
	Space

	// Returns the number of outcomes in the space if finite, and
	// returns -1 if infinite.
	Size() int
}

// A discrete subset of the reals
type DiscreteRealSpace interface {
	DiscreteSpace

	// The real value of an outcome
	F64Value(outcome Outcome) float64
}

// A sample space over boolean outcomes
type BooleanSpace struct{}

// Return the cardinality of the space
func (sp BooleanSpace) Size() int {
	return 2
}

// The real value of an outcome
func (sp BooleanSpace) F64Value(outcome Outcome) float64 {
	if outcome == 1 {
		return 1.0
	} else {
		return 0.0
	}
}

// Return the specified outcome as a boolean
func (sp BooleanSpace) BoolValue(outcome Outcome) bool {
	return outcome == 1
}
