package models

import (
	"database/sql/driver"
	"strconv"
	"strings"

	"github.com/unluckythoughts/manga-reader/utils"
)

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
