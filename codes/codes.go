package codes

import (
	"fmt"
	"strconv"
)

type Code uint32

const (
	// OK is returned on success.
	OK                Code = 0
	Canceled          Code = 1
	DeadlineExceeded  Code = 4
	ResourceExhausted Code = 8
	Aborted           Code = 10
	OutOfRange        Code = 11
	Internal          Code = 13
	Unavailable       Code = 14
	DataLoss          Code = 15
	_maxCode               = 17
)

var strToCode = map[string]Code{
	`"OK"`: OK,
	`"CANCELLED"`:/* [sic] */ Canceled,
	`"DEADLINE_EXCEEDED"`:  DeadlineExceeded,
	`"RESOURCE_EXHAUSTED"`: ResourceExhausted,
	`"ABORTED"`:            Aborted,
	`"OUT_OF_RANGE"`:       OutOfRange,
	`"INTERNAL"`:           Internal,
	`"UNAVAILABLE"`:        Unavailable,
	`"DATA_LOSS"`:          DataLoss,
}

var CodeToStr = map[Code]string{
	OK:                "OK",
	Canceled:          "CANCELLED",
	DeadlineExceeded:  "DEADLINE_EXCEEDED",
	ResourceExhausted: "RESOURCE_EXHAUSTED",
	Aborted:           "ABORTED",
	OutOfRange:        "OUT_OF_RANGE",
	Internal:          "INTERNAL",
	Unavailable:       "UNAVAILABLE",
	DataLoss:          "DATA_LOSS",
}

// UnmarshalJSON unmarshals b into the Code.
func (c *Code) UnmarshalJSON(b []byte) error {
	// From json.Unmarshaler: By convention, to approximate the behavior of
	// Unmarshal itself, Unmarshalers implement UnmarshalJSON([]byte("null")) as
	// a no-op.
	if string(b) == "null" {
		return nil
	}
	if c == nil {
		return fmt.Errorf("nil receiver passed to UnmarshalJSON")
	}

	if ci, err := strconv.ParseUint(string(b), 10, 32); err == nil {
		if ci >= _maxCode {
			return fmt.Errorf("invalid code: %q", ci)
		}

		*c = Code(ci)
		return nil
	}

	if jc, ok := strToCode[string(b)]; ok {
		*c = jc
		return nil
	}
	return fmt.Errorf("invalid code: %q", string(b))
}

func (c Code) Error() string {
	return CodeToStr[c]
}
