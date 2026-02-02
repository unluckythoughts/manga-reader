package utils

import (
	"strconv"
)

func CompareNumbers(a, b string) (int, error) {
	numA, err := strconv.Atoi(a)
	if err != nil {
		return 0, err
	}
	numB, err := strconv.Atoi(b)
	if err != nil {
		return 0, err
	}

	if numA > numB {
		return 1, nil
	} else if numA < numB {
		return -1, nil
	}
	return 0, nil
}
