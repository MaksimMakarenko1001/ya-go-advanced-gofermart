package moneys

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Money struct {
	amount int64 // cents number
}

func New(amount int64) Money {
	return Money{amount: amount}
}

func (m Money) Add(other Money) Money {
	return Money{amount: m.amount + other.amount}
}

func (m Money) String() string {
	return fmt.Sprint(m.amount)
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) MarshalJSON() ([]byte, error) {
	units, cents := m.amount/100, m.amount%100
	if cents < 0 {
		cents *= -1
	}

	res := fmt.Sprintf("%d.%02d", units, cents)
	if cents < 10 {
		res = fmt.Sprintf("%d.0%d", units, cents)
	}

	res = strings.TrimRight(res, "0")
	res = strings.TrimRight(res, ".")
	return []byte(res), nil
}

func (m *Money) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte{'-'}) {
		return fmt.Errorf("negative value")
	}
	value := string(data)

	dotIndex := strings.Index(value, ".")
	if dotIndex == -1 {
		// There's not dot
		units, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		m.amount = units * 100
		return nil
	}

	unitsStr := value[:dotIndex]
	centsStr := value[dotIndex+1:]

	if len(centsStr) == 1 {
		centsStr += "0"
	}

	units, err := strconv.ParseInt(unitsStr, 10, 64)
	if err != nil {
		return err
	}

	cents, err := strconv.ParseInt(centsStr, 10, 64)
	if err != nil {
		return err
	}
	if cents < 0 || cents >= 100 {
		return fmt.Errorf("invalid cents")
	}
	if units < 0 {
		cents *= -1
	}

	m.amount = units*100 + cents
	return nil

}
