package utils

import (
	"golang.org/x/exp/constraints"
)

func MakeRange[T constraints.Integer](start, stop, step T) (values []T) {
	for i := start; i <= stop; i = i + step {
		values = append(values, i)
	}

	return values
}
