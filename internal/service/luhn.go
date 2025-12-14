package service

import (
	"strconv"
)

func IsLuhnValid(orderNumber string) bool {
	length := 0
	for _, r := range orderNumber {
		if _, err := strconv.Atoi(string(r)); err != nil {
			return false
		}
		length += 1
	}

	sum := 0
	for i, j := length-1, 0; i >= 0; i, j = i-1, j+1 {
		digit, _ := strconv.Atoi(string(orderNumber[i]))

		if j%2 == 1 { // Удваиваем каждую вторую цифру с конца
			digit *= 2
			if digit > 9 {
				digit -= 9 // Если результат больше 9, вычитаем 9
			}
		}

		sum += digit
	}

	return sum%10 == 0
}
