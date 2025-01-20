package toolkit

func (t *Tools) sum(ints []int) int {
	var sum int
	for _, num := range ints {
		sum += num
	}
	return sum
}
