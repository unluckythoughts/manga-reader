package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	HUMAN_READABLE_DATE_FORMAT = "HUMAN_READABLE_DATE_FORMAT"
)

func GetString(num float64) string {
	strFloat := fmt.Sprintf("%f", num)
	strFloat = strings.TrimRight(strFloat, "0")
	strFloat = strings.TrimRight(strFloat, ".")

	return strFloat
}

func StrAdd(a string, b int) string {
	numA, err := strconv.Atoi(a)
	if err != nil {
		return a + strconv.Itoa(b)
	}
	return strconv.Itoa(numA + b)
}

func IsInternalLink(s string) bool {
	pattern := regexp.MustCompile("^[/?&]")

	return pattern.MatchString(s)
}

func parseHumanReadableFormat(date string) (string, error) {
	t := strings.Split(strings.TrimSpace(date), " ")

	if len(t) < 3 {
		t = append([]string{"1"}, t...)
	}

	if t[2] != "ago" {
		return date, errors.New("human readable date does not end with 'ago'")
	}

	var err error
	var num int64
	var timeSpan int64

	if t[0] == "a" || t[0] == "an" || t[0] == "1" || t[0] == "few" {
		num = -1
	} else {
		num, err = strconv.ParseInt(t[0], 10, 64)
		if err != nil {
			return date, errors.Wrapf(err, "human readable date could not parse num of quantity %s", t[0])
		}

		num = -1 * num
	}

	switch {
	case strings.HasPrefix(t[1], "second"):
		timeSpan = int64(time.Second)
	case strings.HasPrefix(t[1], "minute"):
		timeSpan = int64(time.Minute)
	case strings.HasPrefix(t[1], "hour"):
		timeSpan = int64(time.Hour)
	case strings.HasPrefix(t[1], "day"):
		timeSpan = int64(24 * time.Hour)
	case strings.HasPrefix(t[1], "week"):
		timeSpan = int64(7 * 24 * time.Hour)
	case strings.HasPrefix(t[1], "month"):
		timeSpan = int64(30 * 24 * time.Hour)
	case strings.HasPrefix(t[1], "year"):
		timeSpan = int64(365 * 24 * time.Hour)
	default:
		return date, errors.Errorf("human readable date could not parse time span %s", t[1])
	}

	newDate := time.Now().Add(time.Duration(num * timeSpan)).Format("2006-01-02")

	return newDate, nil
}

func ParseDate(date string, format string) (string, error) {
	if strings.HasSuffix(date, "ago") || format == HUMAN_READABLE_DATE_FORMAT {
		return parseHumanReadableFormat(date)
	}

	if format == "" {
		return date, nil
	}

	t, err := time.Parse(format, date)
	if err != nil {
		return date, err
	}

	return t.Format("2006-01-02"), nil
}
