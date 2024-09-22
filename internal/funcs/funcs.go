package funcs

import (
	"fmt"
	"strconv"
)

// Seq provides a means to iterate over integer sequences. It returns a
// slice of integers from x to y inclusive. If x is greater than y it
// returns an empty slice. If x is equal to y it returns a slice with
// a single element x. The slice is always in ascending order.
// Use the function in a template like this:
// {{ $seq := seq 1 5 }}
// {{ range $seq }}
//
//	{{ . }}
//
// {{ end }}
// This will output:
// 1
// 2
// 3
// 4
// 5
func Seq(x, y int) []int {
	if x > y {
		return []int{}
	}

	seq := make([]int, y-x+1)
	for i := range seq {
		seq[i] = x + i
	}
	return seq
}

func Type(v any) string {
	return fmt.Sprintf("%T", v)
}

func ToInt(value any) (int, error) {
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		if intValue, err := strconv.Atoi(v); err == nil {
			return intValue, nil
		}
		return 0, fmt.Errorf("cannot convert string to int: %s", v)
	default:
		return 0, fmt.Errorf("expected float64, int, or string, got %T", value)
	}
}

// FuncMap provides a global map of all available built-in funcs.
func FuncMap() map[string]any {
	return map[string]any{
		"seq":   Seq,
		"type":  Type,
		"toInt": ToInt,
	}
}
