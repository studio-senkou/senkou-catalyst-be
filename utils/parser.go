package utils

import "strconv"

func StrToUint(str string) (uint, error) {
	parsedStr, err := strconv.ParseUint(str, 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(parsedStr), nil
}