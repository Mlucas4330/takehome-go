package service

import "unicode"

func ValidateCPF(cpf string) bool {
	d := make([]int, 0, 11)
	for _, r := range cpf {
		if unicode.IsDigit(r) {
			d = append(d, int(r-'0'))
		}
	}
	if len(d) != 11 {
		return false
	}
	allEq := true
	for i := 1; i < 11; i++ {
		if d[i] != d[0] {
			allEq = false
			break
		}
	}
	if allEq {
		return false
	}
	sum := 0
	for i := range 9 {
		sum += d[i] * (10 - i)
	}
	dv1 := (sum * 10) % 11
	if dv1 == 10 {
		dv1 = 0
	}
	if dv1 != d[9] {
		return false
	}
	sum = 0
	for i := range 10 {
		sum += d[i] * (11 - i)
	}
	dv2 := (sum * 10) % 11
	if dv2 == 10 {
		dv2 = 0
	}
	return dv2 == d[10]
}
