package reader

import (
	"errors"
	"strconv"
)

var ErrInvalidNumericalValue = errors.New("reader: invalid numerical value")

func toFloat(val string) (float64, error) {
	var final float64

	if i, err := strconv.Atoi(val); err == nil {
		final = float64(i)
	} else if f, err := strconv.ParseFloat(val, 64); err == nil {
		final = f
	} else {
		return final, ErrInvalidNumericalValue
	}

	return final, nil
}

func dmsToDd(degStr, minStr, secStr string) (float64, error) {
	deg, err := toFloat(degStr)
	if err != nil {
		return 0, err
	}
	min, err := toFloat(degStr)
	if err != nil {
		return 0, err
	}
	sec, err := toFloat(secStr)
	if err != nil {
		return 0, err
	}

	return deg + min/60 + sec/3600, nil
}
