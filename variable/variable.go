package variable

import (
	"github.com/jesand/stats/dist"
)

// A random variable in a factor graph
type RandomVariable interface {

	// Get the variable's current value
	Val() float64

	// Set the variable's current value
	Set(val float64)

	// Ask whether the variable has the same domain and value as another variable
	Equals(other RandomVariable) bool
}

// Create a new continuous random variable
func NewContinuousRV(val float64, space dist.RealSpace) *ContinuousRV {
	return &ContinuousRV{
		val:   val,
		space: space,
	}
}

// A continuous random variable
type ContinuousRV struct {
	val   float64
	space dist.RealSpace
}

func (rv ContinuousRV) Val() float64 {
	return rv.val
}

func (rv *ContinuousRV) Set(val float64) {
	rv.val = val
}

func (rv ContinuousRV) Equals(other RandomVariable) bool {
	crv, ok := other.(*ContinuousRV)
	if !ok {
		return false
	}
	return rv.Space().Equals(crv.Space()) && rv.Val() == crv.Val()
}

func (rv ContinuousRV) Space() dist.RealSpace {
	return rv.space
}

// Create a new discrete random variable
func NewDiscreteRV(val dist.Outcome, space dist.DiscreteRealSpace) *DiscreteRV {
	return &DiscreteRV{
		val:   val,
		space: space,
	}
}

// A discrete random variable
type DiscreteRV struct {
	val   dist.Outcome
	space dist.DiscreteRealSpace
}

func (rv DiscreteRV) Val() float64 {
	return rv.space.F64Value(rv.val)
}

func (rv DiscreteRV) Outcome() dist.Outcome {
	return rv.val
}

func (rv *DiscreteRV) Set(val float64) {
	rv.val = rv.space.Outcome(val)
}

func (rv *DiscreteRV) SetOutcome(val dist.Outcome) {
	rv.val = val
}

func (rv DiscreteRV) Equals(other RandomVariable) bool {
	drv, ok := other.(*DiscreteRV)
	if !ok {
		return false
	}
	return rv.Space().Equals(drv.Space()) && rv.Val() == drv.Val()
}

func (rv DiscreteRV) Space() dist.DiscreteRealSpace {
	return rv.space
}
