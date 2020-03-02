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
