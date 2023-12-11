package maths

// LCM returns the lowest common multiple of all the given numbers
func LCM(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}

	a := numbers[0]
	for _, b := range numbers[1:] {
		a = (a * b) / GCD(a, b)
	}

	return a
}

// GCD returns the greatest common divisor of the given numbers
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
