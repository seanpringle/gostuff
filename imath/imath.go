package imath

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Cycle(i, l int) int {
	if i == -1 {
		i = l - 1
	}
	if i == l {
		i = 0
	}
	return i
}

func If(b bool, n, d int) int {
	if b {
		return n
	}
	return d
}

func IfZero(n, d int) int {
	if n == 0 {
		return d
	}
	return n
}

func IfNeg(n, d int) int {
	if n < 0 {
		return d
	}
	return n
}
