package sarray

import "strings"

func Contains(sarray [3]string, toFind string) bool {
	for _, str := range sarray {
		if strings.EqualFold(str, toFind) {
			return true
		}
	}
	return false
}
