package dist

import (
	"github.com/jesand/stats"
	"math"
)

// An ID for a particular outcome in a space
type Outcome int

// A set of outcomes for some probability measure.
type Space interface {

	// Ask whether the space is the same as some other space
	Equals(other Space) bool
}

// Methods contained by spaces over real values
type RealLikeSpace interface {

	// The infimum (min) value in the space, or negative infinity
	Inf() float64

	// The supremum (max) value in the space, or positive infinity
	Sup() float64
}

// A subset of the reals
type RealSpace interface {
	Space
	RealLikeSpace
}

// The space of reals greater than zero
type positiveRealSpace struct{}

// The canonical instance of positiveRealSpace
var PositiveRealSpace positiveRealSpace

// The infimum (min) value in the space, or negative infinity
func (sp positiveRealSpace) Inf() float64 {
	return 0
}

// The supremum (max) value in the space, or positive infinity
func (sp positiveRealSpace) Sup() float64 {
	return math.Inf(+1)
}

// Ask whether the space is the same as some other space
func (sp positiveRealSpace) Equals(other Space) bool {
	if _, ok := other.(*positiveRealSpace); ok {
		return true
	} else if _, ok := other.(positiveRealSpace); ok {
		return true
	}
	return false
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

// Ask whether the space is the same as some other space
func (sp RealIntervalSpace) Equals(other Space) bool {
	ris, ok := other.(*RealIntervalSpace)
	if !ok {
		ri, ok := other.(RealIntervalSpace)
		if !ok {
			return false
		} else {
			ris = &ri
		}
	}
	return sp.Min == ris.Min && sp.Max == ris.Max
}

// The space of all reals
var AllRealSpace = RealIntervalSpace{Min: math.Inf(-1), Max: math.Inf(+1)}

// The canonical unit interval space
var UnitIntervalSpace = RealIntervalSpace{Min: 0, Max: 1}

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
	RealLikeSpace

	// The real value of an outcome
	F64Value(outcome Outcome) float64

	// The outcome corresponding to a real value
	Outcome(value float64) Outcome
}

// A sample space over boolean outcomes
type booleanSpace struct{}

// The canonical instance of booleanSpace
var BooleanSpace booleanSpace

// The infimum (min) value in the space, or negative infinity
func (sp booleanSpace) Inf() float64 {
	return 0
}

// The supremum (max) value in the space, or positive infinity
func (sp booleanSpace) Sup() float64 {
	return 1
}

// Ask whether the space is the same as some other space
func (sp booleanSpace) Equals(other Space) bool {
	if _, ok := other.(*booleanSpace); ok {
		return true
	} else if _, ok := other.(booleanSpace); ok {
		return true
	}
	return false
}

// Return the cardinality of the space
func (sp booleanSpace) Size() int {
	return 2
}

// The real value of an outcome
func (sp booleanSpace) F64Value(outcome Outcome) float64 {
	if outcome == 0 {
		return 0.0
	} else {
		return 1.0
	}
}

// The outcome corresponding to a real value
func (sp booleanSpace) Outcome(value float64) Outcome {
	if value == 0.0 {
		return 0
	} else {
		return 1
	}
}

// Return the specified outcome as a boolean
func (sp booleanSpace) BoolValue(outcome Outcome) bool {
	return outcome != 0
}

// Return the outcome corresponding to the provided boolean value
func (sp booleanSpace) BoolOutcome(value bool) Outcome {
	if value {
		return 1
	} else {
		return 0
	}
}

// A discrete space over arbitrary objects
type DiscreteObjectSpace struct {

	// The objects which the space is over
	Objects []interface{}
}

// Ask whether the space is the same as some other space
func (sp DiscreteObjectSpace) Equals(other Space) bool {
	var sp2 *DiscreteObjectSpace
	if s, ok := other.(*DiscreteObjectSpace); ok {
		sp2 = s
	} else if s, ok := other.(DiscreteObjectSpace); ok {
		sp2 = &s
	} else {
		return false
	}
	if len(sp.Objects) != len(sp2.Objects) {
		return false
	}
	for i, v := range sp.Objects {
		if v != sp2.Objects[i] {
			return false
		}
	}
	return true
}

// Returns the number of outcomes in the space if finite, and
// returns -1 if infinite.
func (sp DiscreteObjectSpace) Size() int {
	return len(sp.Objects)
}

func (sp DiscreteObjectSpace) Outcome(val interface{}) Outcome {
	for i, v := range sp.Objects {
		if v == val {
			return Outcome(i)
		}
	}
	panic(stats.ErrfValNotInDomain(val))
}

func (sp DiscreteObjectSpace) Value(outcome Outcome) interface{} {
	if int(outcome) < 0 || int(outcome) >= len(sp.Objects) {
		panic(stats.ErrfNotInDomain(int(outcome)))
	}
	return sp.Objects[int(outcome)]
}
