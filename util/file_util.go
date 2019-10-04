package util

import (
	"fmt"
	regexp2 "regexp"
	"strconv"
	"time"
)

func CheckFileName(fileName string) bool {
	now := time.Now()
	regexp := fmt.Sprintf(`(?m)MethodServer\d{1,2}-%s%s%s\d{1,99}-\d{1,99}-log4j.log`,
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
	matched, _ := regexp2.MatchString(regexp, fileName)
	return matched
}

func CheckFileNameOmittingDate(fileName string) bool {
	matched, _ := regexp2.MatchString(fmt.Sprintf(`(?m)MethodServer\d{1,2}-\d{1,99}-\d{1,99}-log4j.log`), fileName)
	return matched
}
