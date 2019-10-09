package util

import (
	proto "dominikw.pl/wnc_plugin/proto"
	"fmt"
	"github.com/google/logger"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type IntegrationFileOperationsTests struct {
	suite.Suite
	correctFilenameForMS  string
	correctFilenameForBMS string
}

func (s *IntegrationFileOperationsTests) BeforeTest(_, testName string) {
	now := time.Now()
	if !strings.Contains(testName, "Today") {
		now = time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	}
	s.correctFilenameForMS = fmt.Sprintf("MethodServer1-%s%s%s319-30967-log4j.log",
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
	s.correctFilenameForBMS = fmt.Sprintf("BackgroundMethodServer1-%s%s%s319-30967-log4j.log",
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
	createFiles(s)
}

func (s *IntegrationFileOperationsTests) AfterTest(_, b string) {
	removeFiles()
}

func TestFindLogFile(t *testing.T) {
	suite.Run(t, new(IntegrationFileOperationsTests))
}

func (suite *IntegrationFileOperationsTests) TestShouldFindBmsStartedToday() {
	msLocation := proto.LogFileLocation{FileLocation: "test", LogType: proto.LogFileLocation_BACKGROUND_METHOD_SERVER}
	s, _ := FindLogFile(&msLocation)
	suite.Equal(suite.correctFilenameForBMS, s, "should be equal")
}

func (suite *IntegrationFileOperationsTests) TestShouldFindMsStartedToday() {
	msLocation := proto.LogFileLocation{FileLocation: "test", LogType: proto.LogFileLocation_METHOD_SERVER}
	s, _ := FindLogFile(&msLocation)
	suite.Equal(suite.correctFilenameForMS, s, "should be equal")
}

func (suite *IntegrationFileOperationsTests) TestShouldFindBmsStartedBefore() {
	msLocation := proto.LogFileLocation{FileLocation: "test", LogType: proto.LogFileLocation_BACKGROUND_METHOD_SERVER}
	s, _ := FindLogFile(&msLocation)
	suite.Equal(suite.correctFilenameForBMS, s, "should be equal")
}

func (suite *IntegrationFileOperationsTests) TestShouldFindMsStartedBefore() {
	msLocation := proto.LogFileLocation{FileLocation: "test", LogType: proto.LogFileLocation_METHOD_SERVER}
	s, _ := FindLogFile(&msLocation)
	suite.Equal(suite.correctFilenameForMS, s, "should be equal")
}

func (suite *IntegrationFileOperationsTests) TestShouldFailWithIncorrectDirectory() {
	msLocation := proto.LogFileLocation{FileLocation: "incorrect", LogType: proto.LogFileLocation_METHOD_SERVER}
	s, e := FindLogFile(&msLocation)
	suite.Assertions.Error(e, "%s Should be an incorrect directory name", msLocation.FileLocation)
	suite.Assertions.Empty(s, "%s Should be empty", s)
}

func (suite *IntegrationFileOperationsTests) TestShouldFailWithNoFilesFound() {
	infos, _ := ioutil.ReadDir("test")
	for _, info := range infos {
		_ = os.Remove("test/" + info.Name())
	}

	msLocation := proto.LogFileLocation{FileLocation: "test", LogType: proto.LogFileLocation_METHOD_SERVER}
	s, e := FindLogFile(&msLocation)
	suite.Assertions.EqualError(e, "Log file not found for type: METHOD_SERVER!",
		`Should not found any files`, msLocation.FileLocation)
	suite.Assertions.Empty(s, "%s Should be empty", s)
}

func removeFiles() {
	e := os.RemoveAll("test")
	if e != nil {
		logger.Fatal(e)
	}
}

func createFiles(s *IntegrationFileOperationsTests) {
	_ = os.Mkdir("test", 777)
	createFile("test/" + s.correctFilenameForMS)
	createFile("test/" + s.correctFilenameForBMS)
}

func createFile(name string) {
	f, e := os.OpenFile(name, os.O_CREATE, 777)
	if e != nil {
		logger.Fatal(e)
	}
	defer f.Close()
}
