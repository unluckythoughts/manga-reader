package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/exp/constraints"
)

const (
	HUMAN_READABLE_DATE_FORMAT = "HUMAN_READABLE_DATE_FORMAT"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func GetInt(val string) int {
	pattern := regexp.MustCompile("[^0-9]")
	strNum := pattern.ReplaceAllString(val, "")

	return MustParse[int](strNum)
}

func MustParse[T Number](s string) T {
	v, _ := Parse[T](s)
	return v
}

func Parse[T Number](s string) (T, error) {
	var ret T
	switch any(ret).(type) {
	case int:
		v, err := strconv.ParseInt(s, 10, 0)
		return T(v), err
	case int8:
		v, err := strconv.ParseInt(s, 10, 8)
		return T(v), err
	case int16:
		v, err := strconv.ParseInt(s, 10, 16)
		return T(v), err
	case int32:
		v, err := strconv.ParseInt(s, 10, 32)
		return T(v), err
	case int64:
		v, err := strconv.ParseInt(s, 10, 64)
		return T(v), err
	case uint:
		v, err := strconv.ParseInt(s, 10, 0)
		return T(v), err
	case uint8:
		v, err := strconv.ParseInt(s, 10, 8)
		return T(v), err
	case uint16:
		v, err := strconv.ParseInt(s, 10, 16)
		return T(v), err
	case uint32:
		v, err := strconv.ParseInt(s, 10, 32)
		return T(v), err
	case uint64:
		v, err := strconv.ParseInt(s, 10, 64)
		return T(v), err
	case float32:
		v, err := strconv.ParseFloat(s, 32)
		return T(v), err
	case float64:
		v, err := strconv.ParseFloat(s, 64)
		return T(v), err
	}

	return ret, errors.Errorf("unexptect type %T", ret)
}

func GetString(val interface{}) string {
	switch v := val.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case int8:
		return fmt.Sprintf("%d", v)
	case int16:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case uint:
		return fmt.Sprintf("%d", v)
	case uint8:
		return fmt.Sprintf("%d", v)
	case uint16:
		return fmt.Sprintf("%d", v)
	case uint32:
		return fmt.Sprintf("%d", v)
	case uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		strFloat := fmt.Sprintf("%.6f", v)
		strFloat = strings.TrimRight(strFloat, "0")
		return strings.TrimRight(strFloat, ".")
	case float64:
		strFloat := fmt.Sprintf("%.6f", v)
		strFloat = strings.TrimRight(strFloat, "0")
		return strings.TrimRight(strFloat, ".")
	}

	return ""
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

	if format == "" || date == "" {
		return date, nil
	}

	t, err := time.Parse(format, date)
	if err != nil {
		return date, err
	}

	return t.Format("2006-01-02"), nil
}
