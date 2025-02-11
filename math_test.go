package toolkit

import "testing"

func Test_Sum(t *testing.T) {
	var tools Tools
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"Small ints", []int{1, 2, 3, 4, 5}, 15},
		{"Mixed ints", []int{1, 200, 3, 499, 5}, 708},
		{"Small negative ints", []int{-1, -2, -3, 4, -5}, -7},
		{"Mixed negative ints", []int{1, -200, 3, 499, 5}, 308},
		{"Big ints", []int{1111, 200, 3333, 499, 5555}, 10698},
		{"Big mixed ints", []int{-1522, 2777733333, -3242411, 4432224, -1117895}, 2777803729},
		{"Negative big ints", []int{-1111, -2400, -3333, -499, -5555}, -12898},
		{"Zero", []int{0}, 0},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := tools.Sum(entry.nums)

			if result != entry.expected {
				t.Errorf("expected %d, received %d", entry.expected, result)
			}

		})
	}
}
