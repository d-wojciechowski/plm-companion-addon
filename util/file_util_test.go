package util

import (
	proto "dominikw.pl/wnc_plugin/proto"
	"fmt"
	"github.com/google/logger"
	"os"
	"strconv"
	"testing"
	"time"
)

var correctFilenameForMS string
var correctFilenameForBMS string

func init() {
	now := time.Now()
	correctFilenameForMS = fmt.Sprintf("MethodServer1-%s%s%s319-30967-log4j.log",
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
	correctFilenameForBMS = fmt.Sprintf("BackgroundMethodServer1-%s%s%s319-30967-log4j.log",
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
	createFiles()
}

func TestCheckFileName(t *testing.T) {
	t.Run("TestCheckFileNameShouldWork", testCheckFileNameShouldWork)
	t.Run("testCheckFileNameShouldFail", testCheckFileNameShouldFail)
}

func TestFindLogFile(t *testing.T) {

	t.Run("shouldFindMsStartedToday", shouldFindMsStartedToday)
	t.Run("shouldFindBmsStartedToday", shouldFindBmsStartedToday)

	removeFiles()
}

func shouldFindBmsStartedToday(t *testing.T) {
	msLocation := proto.LogFileLocation{FileLocation: "test", LogType: proto.LogFileLocation_BACKGROUND_METHOD_SERVER}
	s, _ := FindLogFile(&msLocation)
	if s == correctFilenameForBMS {
		t.Logf("test OK!")
	} else {
		t.Errorf("Test failed!")
	}
}

func shouldFindMsStartedToday(t *testing.T) {
	msLocation := proto.LogFileLocation{FileLocation: "test", LogType: proto.LogFileLocation_METHOD_SERVER}
	s, _ := FindLogFile(&msLocation)
	if s == correctFilenameForMS {
		t.Logf("test OK!")
	} else {
		t.Errorf("Test failed!")
	}
}

func removeFiles() {
	e := os.RemoveAll("test")
	if e != nil {
		logger.Fatal(e)
	}
}

func createFiles() {
	os.Mkdir("test", 777)
	createFile("test/" + correctFilenameForMS)
	createFile("test/" + correctFilenameForBMS)
}

func createFile(name string) {
	f, e := os.OpenFile(name, os.O_CREATE, 777)
	if e != nil {
		logger.Fatal(e)
	}
	defer f.Close()
}

func testCheckFileNameShouldWork(t *testing.T) {
	if checkFileName(correctFilenameForMS, proto.LogFileLocation_METHOD_SERVER) {
		t.Logf("%s: regexp OK!", t.Name())
	} else {
		t.Errorf("%s should work for %s", t.Name(), correctFilenameForMS)
	}
}

func testCheckFileNameShouldFail(t *testing.T) {
	name := "random test string"
	if !checkFileName(name, proto.LogFileLocation_METHOD_SERVER) {
		t.Logf("%s: regexp OK!", t.Name())
	} else {
		t.Errorf(`%s should not work for "%s"`, t.Name(), name)
	}
}
