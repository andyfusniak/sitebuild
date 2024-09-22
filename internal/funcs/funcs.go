package funcs

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

// FuncMap provides a global map of all available built-in funcs.
func FuncMap() map[string]any {
	return map[string]any{
		"seq": Seq,
	}
}
