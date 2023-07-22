package models

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/unluckythoughts/manga-reader/utils"
)

type StrTimeStamp time.Time

func (sts StrTimeStamp) MarshalJSON() ([]byte, error) {
	ts := "\"" + time.Time(sts).Format(time.DateTime) + "\""
	return []byte(ts), nil
}

func (sts *StrTimeStamp) Scan(value interface{}) error {
	switch value := value.(type) {
	case time.Time:
		*sts = StrTimeStamp(value)
		return nil
	case string:
		numPattern := regexp.MustCompile(`^[0-9]+$`)
		if numPattern.MatchString(value) {
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			*sts = StrTimeStamp(time.Unix(i, 0))
			return nil
		}
		ts, err := time.Parse(time.RFC3339, value)
		*sts = StrTimeStamp(ts)
		return err
	}

	return fmt.Errorf("could not parse StrTimeStamp value from timestamp: %v of type: %T", value, value)
}

func (sts StrTimeStamp) Value() (driver.Value, error) {
	ts := time.Time(sts)

	return ts.Format(time.RFC3339), nil
}

type StrFloat string

func (sf *StrFloat) UnmarshalJSON(data []byte) error {
	*sf = StrFloat(data)

	return nil
}

type StrFloatList [2]float64

func (sfl *StrFloatList) Scan(value interface{}) error {
	l := StrFloatList{0, 0}
	data, ok := value.(string)
	if !ok || data == "" {
		return nil
	}

	sl := strings.Split(data, ",")
	for i, s := range sl {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		l[i] = v
	}

	*sfl = l
	return nil
}

func (sfl StrFloatList) Value() (driver.Value, error) {
	if len(sfl) == 0 {
		return "", nil
	}

	l := []string{}
	for _, v := range sfl {
		l = append(l, utils.GetString(v))
	}

	return strings.Join(l, ","), nil
}

type StrIntList [2]int

func (sil *StrIntList) Scan(value interface{}) error {
	l := StrIntList{0, 0}
	data, ok := value.(string)
	if !ok || data == "" {
		return nil
	}

	sl := strings.Split(data, ",")
	for i, s := range sl {
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		l[i] = v
	}

	*sil = l
	return nil
}

func (sil StrIntList) Value() (driver.Value, error) {
	if len(sil) == 0 {
		return "", nil
	}

	l := []string{}
	for _, v := range sil {
		l = append(l, strconv.Itoa(v))
	}

	return strings.Join(l, ","), nil
}
