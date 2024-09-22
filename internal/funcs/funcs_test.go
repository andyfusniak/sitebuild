package funcs_test

import (
	"slices"
	"testing"

	"github.com/andyfusniak/sitebuild/internal/funcs"
)

func TestSeq(t *testing.T) {
	tests := []struct {
		x, y     int
		expected []int
	}{
		{x: 1, y: 5, expected: []int{1, 2, 3, 4, 5}},
		{x: 4, y: 9, expected: []int{4, 5, 6, 7, 8, 9}},
		{x: 5, y: 1, expected: []int{}},
		{x: 1, y: 1, expected: []int{1}},
	}

	for i, tt := range tests {
		seq := funcs.Seq(tt.x, tt.y)
		if !slices.Equal(seq, tt.expected) {
			t.Errorf("Seq(%d, %d) got %v, expected %v for test %d",
				tt.x, tt.y, seq, tt.expected, i)
		}
	}
}
