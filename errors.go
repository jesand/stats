package stats

import (
	"fmt"
)

const (
	ErrNotNormalized Error = "The distribution was not normalized properly"
	ErrZeroProb      Error = "The distribution has zero total probability"
)

func ErrfNotInDomain(outcome int) Error {
	return Errorf("Outcome %d not in the sample space", outcome)
}

func ErrfInvalidProb(prob float64) Error {
	return Errorf("Invalid discrete probability %f", prob)
}

func ErrfFactorVarNum(numVars, numParams, numAdj int) Error {
	return Errorf("Factor expected %d variables and %d parameters, but has %d adjacent",
		numVars, numParams, numAdj)
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
