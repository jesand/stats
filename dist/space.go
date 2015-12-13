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

// Create a new RealIntervalSpace with the specified bounds
func NewRealIntervalSpace(min, max float64) *RealIntervalSpace {
	return &RealIntervalSpace{
		Min: min,
		Max: max,
	}
}

// A subset of the reals on a continuous closed interval
type RealIntervalSpace struct {
	Min, Max float64
}

// The infimum (min) value in the space, or negative infinity
func (space RealIntervalSpace) Inf() float64 {
	return space.Min
}

// The supremum (max) value in the space, or positive infinity
func (space RealIntervalSpace) Sup() float64 {
	return space.Max
}

// The unit interval
func NewUnitIntervalSpace() RealIntervalSpace {
	return RealIntervalSpace{Min: 0, Max: 1}
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
	RealSpace

	// The real value of an outcome
	F64Value(outcome Outcome) float64

	// The outcome corresponding to a real value
	Outcome(value float64) Outcome
}

// A sample space over boolean outcomes
type BooleanSpace struct{}

// The infimum (min) value in the space, or negative infinity
func (sp BooleanSpace) Inf() float64 {
	return 0
}

// The supremum (max) value in the space, or positive infinity
func (sp BooleanSpace) Sup() float64 {
	return 1
}

// Return the cardinality of the space
func (sp BooleanSpace) Size() int {
	return 2
}

// The real value of an outcome
func (sp BooleanSpace) F64Value(outcome Outcome) float64 {
	if outcome == 0 {
		return 0.0
	} else {
		return 1.0
	}
}

// The outcome corresponding to a real value
func (sp BooleanSpace) Outcome(value float64) Outcome {
	if value == 0.0 {
		return 0
	} else {
		return 1
	}
}

// Return the specified outcome as a boolean
func (sp BooleanSpace) BoolValue(outcome Outcome) bool {
	return outcome != 0
}
