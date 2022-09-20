package flextime

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
)

type CombinedFlextime struct {
	parsers   []*Flextime
	numParser func(int64) time.Time
}

func NewCombined(parsers []*Flextime, numParser func(int64) time.Time) *CombinedFlextime {
	return &CombinedFlextime{
		parsers:   parsers,
		numParser: numParser,
	}
}

func (c *CombinedFlextime) Parse(v any) (time.Time, error) {
	return c.parse(v, false, nil)
}

func (c *CombinedFlextime) ParseInLocation(v any, loc *time.Location) (time.Time, error) {
	return c.parse(v, true, loc)
}

func (c *CombinedFlextime) parse(v any, inLoc bool, loc *time.Location) (time.Time, error) {
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if c.numParser != nil {
			return c.numParser(rv.Int()), nil
		}
		return time.Time{}, ErrEmptyNumParser
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if rv.Uint() > math.MaxInt64 {
			return time.Time{}, &ValueOutOfRangeError{Value: rv.Uint()}
		}
		if c.numParser != nil {
			return c.numParser(int64(rv.Uint())), nil
		}
		return time.Time{}, ErrEmptyNumParser
	case reflect.Float32, reflect.Float64:
		if c.numParser != nil {
			// let's simply ignore fraction of number.
			return c.numParser(int64(rv.Float())), nil
		}
		return time.Time{}, ErrEmptyNumParser
	case reflect.String:
		value := rv.String()
		var lastErr error
		for _, f := range c.parsers {
			var parsed time.Time
			var err error
			if inLoc {
				parsed, err = f.ParseInLocation(value, loc)
			} else {
				parsed, err = f.Parse(value)
			}
			if err != nil {
				lastErr = err
			} else {
				return parsed, nil
			}
		}
		return time.Time{}, lastErr
	case reflect.Slice:
		if bs, ok := v.([]byte); ok {
			var jsonVar any
			err := json.Unmarshal(bs, &jsonVar)
			if err != nil {
				return time.Time{}, &UnmarshalError{Err: err}
			}
			switch x := jsonVar.(type) {
			case float64:
				if c.numParser != nil {
					return c.numParser(int64(x)), nil
				}
				return time.Time{}, ErrEmptyNumParser
			case string:
				var lastErr error
				for _, f := range c.parsers {
					var parsed time.Time
					var err error
					if inLoc {
						parsed, err = f.ParseInLocation(x, loc)
					} else {
						parsed, err = f.Parse(x)
					}
					if err != nil {
						lastErr = err
					} else {
						return parsed, nil
					}
				}
				return time.Time{}, lastErr
			}
		}
	}
	return time.Time{}, &UnsupportedTypeError{Typ: rv.Kind()}
}

var ErrEmptyNumParser = errors.New("empty num parser")

type ValueOutOfRangeError struct {
	Value uint64
}

func (e *ValueOutOfRangeError) Error() string {
	return fmt.Sprintf(
		"value out of range: value must be less than max of int64 but is %d",
		e.Value,
	)
}

type UnsupportedTypeError struct {
	Typ reflect.Kind
}

func (e *UnsupportedTypeError) Error() string {
	return fmt.Sprintf(
		"unsupported type: value must be one of Int, Uint, Float variant, String"+
			" or []byte which represents a valid json value and that can be unmarshalled into"+
			" string or float64 but is %s",
		e.Typ.String())
}

type UnmarshalError struct {
	Err error
}

func (e UnmarshalError) Error() string {
	return fmt.Sprintf("unmarshal failed: %v", e.Err)
}
