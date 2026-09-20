package plist

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ParseJSONArray validates plist scalar kinds and rejects values that would be dropped.
func ParseJSONArray(raw string) ([]any, error) {
	if !json.Valid([]byte(raw)) {
		return nil, fmt.Errorf("array_json must contain a valid JSON array")
	}
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	var array []any
	if err := decoder.Decode(&array); err != nil || array == nil {
		return nil, fmt.Errorf("array_json must contain a JSON array")
	}
	value, err := parseArrayValue(array)
	if err != nil {
		return nil, err
	}
	return value.([]any), nil
}

func parseArrayValue(value any) (any, error) {
	switch v := value.(type) {
	case json.Number:
		if strings.ContainsAny(string(v), ".eE") {
			n, err := v.Float64()
			if err != nil || math.IsInf(n, 0) || math.IsNaN(n) {
				return nil, fmt.Errorf("array_json real %q is outside the supported range", v)
			}
			mantissa := strings.SplitN(strings.ToLower(string(v)), "e", 2)[0]
			if n == 0 && strings.ContainsAny(mantissa, "123456789") {
				return nil, fmt.Errorf("array_json real %q underflows the supported range", v)
			}
			return n, nil
		}
		if n, err := v.Int64(); err == nil {
			return n, nil
		}
		if n, err := strconv.ParseUint(string(v), 10, 64); err == nil {
			return n, nil
		}
		return nil, fmt.Errorf("array_json integer %q is outside the supported 64-bit range", v)
	case []any:
		result := make([]any, len(v))
		for i, child := range v {
			var err error
			result[i], err = parseArrayValue(child)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case map[string]any:
		result := make(map[string]any, len(v))
		for key, child := range v {
			var err error
			result[key], err = parseArrayValue(child)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case string, bool:
		return v, nil
	default:
		return nil, fmt.Errorf("array_json contains unsupported value %T; plist does not support null", value)
	}
}

// EncodeJSONArray preserves array order, duplicates, and plist integer/real types.
func EncodeJSONArray(array []any) (string, error) {
	value, err := arrayJSONValue(array)
	if err != nil {
		return "", err
	}
	encoded, err := json.Marshal(value)
	return string(encoded), err
}

func arrayJSONValue(value any) (any, error) {
	switch v := value.(type) {
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, fmt.Errorf("plist array contains a non-finite real")
		}
		raw := strconv.FormatFloat(v, 'g', -1, 64)
		// JSON's default float encoder drops .0, which would turn a plist real into
		// an integer on the next construction and change the set element identity.
		if !strings.ContainsAny(raw, ".eE") {
			raw += ".0"
		}
		return json.Number(raw), nil
	case float32:
		return arrayJSONValue(float64(v))
	case string, bool, int, int64, uint64:
		return v, nil
	case []any:
		result := make([]any, len(v))
		for i, child := range v {
			var err error
			result[i], err = arrayJSONValue(child)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case map[string]any:
		result := make(map[string]any, len(v))
		for key, child := range v {
			var err error
			result[key], err = arrayJSONValue(child)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("plist array contains unsupported JSON value type %T", value)
	}
}
