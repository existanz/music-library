package models

import (
	"fmt"
	"time"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format("2006-01-02") + `"`), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	fmt.Println(string(b))
	*d = NewDateFromString(string(b[1 : len(b)-1]))
	return nil
}

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d Date) Scanner() string {
	return time.Time(d).Format("2006-01-02")
}

func (d Date) Value() (string, error) {
	return time.Time(d).Format("2006-01-02"), nil
}

func NewDateFromString(date string) Date {
	res, err := time.Parse("date", date)
	if err != nil {
		res = time.Now()
	}
	return Date(res)
}
