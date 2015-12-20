package stats

import (
	"fmt"
)

const (
	ErrNotNormalized  Error = "The distribution was not normalized properly"
	ErrZeroProb       Error = "The distribution has zero total probability"
	ErrGraphNotTree   Error = "The factor graph is not a tree"
	ErrDiscreteOnly   Error = "This process currently only supports discrete random variables"
	ErrContinuousOnly Error = "This process currently only supports continuous random variables"
	ErrBernoulliOnly  Error = "This process only supports Bernoulli random variables"
)

func ErrfNotInDomain(outcome int) Error {
	return Errorf("Outcome %d not in the sample space", outcome)
}

func ErrfValNotInDomain(value interface{}) Error {
	return Errorf("Value %v not in the sample space", value)
}

func ErrfInvalidProb(prob float64) Error {
	return Errorf("Invalid probability %f", prob)
}

func ErrfFactorVarNum(numVars, numParams, numAdj int) Error {
	return Errorf("Factor expected %d variable(s) and %d parameter(s), but has %d adjacent",
		numVars, numParams, numAdj)
}

func ErrfUnsupportedDist(dist interface{}) Error {
	return Errorf("Unsupported distribution type %T", dist)
}

// An error message
type Error string

// Return the error message
func (err Error) Error() string {
	return string(err)
}

// Format an error string
func Errorf(message string, args ...interface{}) Error {
	return Error(fmt.Sprintf(message, args...))
}
