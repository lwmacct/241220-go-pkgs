package mtos

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type this struct{}

func New() *this {
	return &this{}
}

// String converts any basic type to its string representation.
// For float types, it formats the number to two decimal places.
func (t *this) String(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32:
		// Format float32 with two decimal places
		return fmt.Sprintf("%.2f", val)
	case float64:
		// Format float64 with two decimal places
		return fmt.Sprintf("%.2f", val)
	case bool:
		if val {
			return "1"
		}
		return "0"
	default:
		// Attempt JSON marshalling for other types
		data, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		return string(data)
	}
}

// Int converts a string or numeric type to int.
func (t *this) Int(v interface{}) int {
	switch val := v.(type) {
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
		return 0
	case int:
		return val
	case int8, int16, int32, int64:
		return int(t.Int64(val))
	case uint, uint8, uint16, uint32, uint64:
		return int(t.Uint64(val))
	case float32, float64:
		return int(t.Float64(val))
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// Int64 converts a string or numeric type to int64.
func (t *this) Int64(v interface{}) int64 {
	switch val := v.(type) {
	case string:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
		return 0
	case int:
		return int64(val)
	case int8:
		return int64(val)
	case int16:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	case uint, uint8, uint16, uint32, uint64:
		return int64(t.Uint64(val))
	case float32:
		return int64(val)
	case float64:
		return int64(val)
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// Float64 converts a string or numeric type to float64.
func (t *this) Float64(v interface{}) float64 {
	switch val := v.(type) {
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
		return 0
	case int, int8, int16, int32, int64:
		return float64(t.Int64(val))
	case uint, uint8, uint16, uint32, uint64:
		return float64(t.Uint64(val))
	case float32:
		return float64(val)
	case float64:
		return val
	case bool:
		if val {
			return 1.0
		}
		return 0.0
	default:
		return 0
	}
}

// Bool converts a string or numeric type to bool.
func (t *this) Bool(v interface{}) bool {
	switch val := v.(type) {
	case string:
		return val == "1" || val == "true" || val == "True"
	case int, int8, int16, int32, int64:
		return t.Int64(val) == 1
	case uint, uint8, uint16, uint32, uint64:
		return t.Uint64(val) == 1
	case float32, float64:
		return t.Float64(val) == 1.0
	case bool:
		return val
	default:
		return false
	}
}

// Uint converts a string or numeric type to uint.
func (t *this) Uint(v interface{}) uint {
	switch val := v.(type) {
	case string:
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			return uint(u)
		}
		return 0
	case int:
		return uint(val)
	case int8, int16, int32, int64:
		return uint(t.Int64(val))
	case uint, uint8, uint16, uint32, uint64:
		switch v := val.(type) {
		case uint:
			return uint(v)
		case uint8:
			return uint(v)
		case uint16:
			return uint(v)
		case uint32:
			return uint(v)
		case uint64:
			return uint(v)
		}
	case float32, float64:
		return uint(t.Float64(val))
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		return 0
	}
	return 0
}

// Uint64 converts a string or numeric type to uint64.
func (t *this) Uint64(v interface{}) uint64 {
	switch val := v.(type) {
	case string:
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			return u
		}
		return 0
	case int, int8, int16, int32, int64:
		return uint64(t.Int64(val))
	case uint, uint8, uint16, uint32, uint64:
		switch v := val.(type) {
		case uint:
			return uint64(v)
		case uint8:
			return uint64(v)
		case uint16:
			return uint64(v)
		case uint32:
			return uint64(v)
		case uint64:
			return v
		}
	case float32, float64:
		return uint64(t.Float64(val))
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		return 0
	}
	return 0
}

// Float32 converts to float32.
func (t *this) Float32(v interface{}) float32 {
	return float32(t.Float64(v))
}

// Json converts any type to its JSON string representation.
// It pretty-prints the JSON.
func (t *this) Json(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}

// Add more "to" methods as needed, for example:

// Uint32 converts to uint32.
func (t *this) Uint32(v interface{}) uint32 {
	return uint32(t.Uint64(v))
}