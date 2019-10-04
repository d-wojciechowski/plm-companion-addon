package util

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestCheckFileNameShouldWork(t *testing.T) {
	now := time.Now()

	name := fmt.Sprintf("MethodServer1-%s%s%s319-30967-log4j.log",
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
	if CheckFileName(name) {
		t.Logf("%s: regexp OK!", t.Name())
	} else {
		t.Errorf("%s should work for %s", t.Name(), name)
	}
}

func TestCheckFileNameShouldFail(t *testing.T) {
	name := "random test string"
	if !CheckFileName(name) {
		t.Logf("%s: regexp OK!", t.Name())
	} else {
		t.Errorf(`%s should not work for "%s"`, t.Name(), name)
	}
}
