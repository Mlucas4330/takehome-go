package validator

import (
	"strconv"
	"strings"
)

func ValidateCPF(cpf string) bool {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")

	if len(cpf) != 11 {
		return false
	}

	if _, err := strconv.ParseUint(cpf, 10, 64); err != nil {
		return false
	}

	allSame := true
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != cpf[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	sum := 0
	for i := 0; i < 9; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * (10 - i)
	}
	digit1 := (sum * 10) % 11
	if digit1 == 10 {
		digit1 = 0
	}

	sum = 0
	for i := 0; i < 10; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * (11 - i)
	}
	digit2 := (sum * 10) % 11
	if digit2 == 10 {
		digit2 = 0
	}

	d1, _ := strconv.Atoi(string(cpf[9]))
	d2, _ := strconv.Atoi(string(cpf[10]))

	return d1 == digit1 && d2 == digit2
}
