package util

import (
	"github.com/google/logger"
)

func IsEmpty(s string) bool {
	return &s == nil || s == ""
}

func PanicOnError(e error) {
	if e != nil {
		logger.Error(e)
		panic(e)
	}
}

func PanicWrapper(i interface{}, e error) interface{} {
	PanicOnError(e)
	return i
}

func Filter(arr []string, cond func(string) bool) []string {
	var result []string
	for i := range arr {
		if cond(arr[i]) {
			result = append(result, arr[i])
		}
	}
	return result
}
