//go:build !solution

package speller

import (
	"math"
	"strings"
)

var dictionary = map[int64]string{
	1:          "one",
	2:          "two",
	3:          "three",
	4:          "four",
	5:          "five",
	6:          "six",
	7:          "seven",
	8:          "eight",
	9:          "nine",
	10:         "ten",
	11:         "eleven",
	12:         "twelve",
	13:         "thirteen",
	14:         "fourteen",
	15:         "fifteen",
	16:         "sixteen",
	17:         "seventeen",
	18:         "eighteen",
	19:         "nineteen",
	20:         "twenty",
	30:         "thirty",
	40:         "forty",
	50:         "fifty",
	60:         "sixty",
	70:         "seventy",
	80:         "eighty",
	90:         "ninety",
	100:        "hundred",
	1000:       "thousand",
	1000000:    "million",
	1000000000: "billion"}

func numberToTriad(n int64) string {
	n = int64(math.Abs(float64(n)))
	var hundreds, dozens, units int64
	var result string

	hundreds = n / 100
	dozens = (n % 100) / 10
	units = n % 10

	if hundreds != 0 {
		result += dictionary[hundreds] + " " + dictionary[100] + " "
	}

	if dozens == 0 {
		result += dictionary[units]
	} else {
		if dozens == 1 {
			result += dictionary[dozens*10+units]
		}

		if dozens > 1 && (dozens*10+units)%10 != 0 {
			result += dictionary[dozens*10] + "-" + dictionary[units]
		}

		if dozens > 1 && (dozens*10+units)%10 == 0 {
			result += dictionary[dozens*10]
		}
	}

	return result
}

func Spell(n int64) string {
	var result string

	if n == 0 {
		return "zero"
	}

	if n < 0 {
		result += "minus" + " "
	}

	n = int64(math.Abs(float64(n)))
	var billions, millions, thousands, units int64

	billions = n / 1000000000
	millions = (n % 1000000000) / 1000000
	thousands = ((n % 1000000000) % 1000000) / 1000
	units = ((n % 1000000000) % 1000000) % 1000

	if billions != 0 {
		result += numberToTriad(billions) + " " + dictionary[1000000000] + " "
	}

	if millions != 0 {
		result += numberToTriad(millions) + " " + dictionary[1000000] + " "
	}

	if thousands != 0 {
		result += numberToTriad(thousands) + " " + dictionary[1000] + " "
	}

	if units != 0 {
		result += numberToTriad(units)
	}

	if strings.HasSuffix(result, " ") {
		result = strings.TrimRight(result, " ")
	}

	return result
}
