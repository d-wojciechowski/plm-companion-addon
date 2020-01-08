package util

import (
	proto "dominikw.pl/wnc_plugin/proto"
	"fmt"
	"github.com/google/logger"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

type IntegrationFileOperationsTests struct {
	suite.Suite
	testDirName      string
	correctFileNames map[proto.LogFileLocation_Source]string
}

func (suite *IntegrationFileOperationsTests) BeforeTest(_, testName string) {
	suite.correctFileNames = make(map[proto.LogFileLocation_Source]string)
	now := time.Now()
	if !strings.Contains(testName, "Today") {
		now = time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	}
	suite.correctFileNames[proto.LogFileLocation_METHOD_SERVER] = fmt.Sprintf(
		"MethodServer1-%s%s%s319-30967-log4j.log",
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
	suite.correctFileNames[proto.LogFileLocation_BACKGROUND_METHOD_SERVER] = fmt.Sprintf(
		"BackgroundMethodServer1-%s%s%s319-30967-log4j.log",
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
	suite.correctFileNames[proto.LogFileLocation_CUSTOM] = "customTestFile.log"
	createFiles(suite)
}

func (suite *IntegrationFileOperationsTests) AfterTest(_, _ string) {
	removeFiles(suite)
}

func TestFindLogFile(t *testing.T) {
	suite.Run(t, new(IntegrationFileOperationsTests))
}

func (suite *IntegrationFileOperationsTests) TestShouldFindBmsStartedToday() {
	msLocation := proto.LogFileLocation{FileLocation: suite.testDirName, LogType: proto.LogFileLocation_BACKGROUND_METHOD_SERVER}
	s, e := FindLogFile(&msLocation)
	if suite.NoError(e, "File should be found!") {
		suite.Equal(suite.correctFileNames[proto.LogFileLocation_BACKGROUND_METHOD_SERVER], s, "should be equal")
	}
}

func (suite *IntegrationFileOperationsTests) TestShouldFindMsStartedToday() {
	msLocation := proto.LogFileLocation{FileLocation: suite.testDirName, LogType: proto.LogFileLocation_METHOD_SERVER}
	s, e := FindLogFile(&msLocation)
	if suite.NoError(e, "File should be found!") {
		suite.Equal(suite.correctFileNames[proto.LogFileLocation_METHOD_SERVER], s, "should be equal")
	}
}

func (suite *IntegrationFileOperationsTests) TestShouldFindBmsStartedBefore() {
	msLocation := proto.LogFileLocation{FileLocation: suite.testDirName, LogType: proto.LogFileLocation_BACKGROUND_METHOD_SERVER}
	s, e := FindLogFile(&msLocation)
	if suite.NoError(e, "File should be found!") {
		suite.Equal(suite.correctFileNames[proto.LogFileLocation_BACKGROUND_METHOD_SERVER], s, "should be equal")
	}
}

func (suite *IntegrationFileOperationsTests) TestShouldFindMsStartedBefore() {
	msLocation := proto.LogFileLocation{FileLocation: suite.testDirName, LogType: proto.LogFileLocation_METHOD_SERVER}
	s, e := FindLogFile(&msLocation)
	if suite.NoError(e, "File should be found!") {
		suite.Equal(suite.correctFileNames[proto.LogFileLocation_METHOD_SERVER], s, "should be equal")
	}
}

func (suite *IntegrationFileOperationsTests) TestShouldFindExistingCustomLog() {
	testFileFullPath := filepath.Join(suite.testDirName, suite.correctFileNames[proto.LogFileLocation_CUSTOM])
	msLocation := proto.LogFileLocation{
		FileLocation: testFileFullPath,
		LogType:      proto.LogFileLocation_CUSTOM}
	s, e := FindLogFile(&msLocation)
	if suite.NoError(e, "File should be found!") {
		suite.Equal(testFileFullPath, s,
			"File name should be %s not %s",
			suite.correctFileNames[proto.LogFileLocation_CUSTOM], s)
	}
}

func (suite *IntegrationFileOperationsTests) TestShouldNotFindExistingCustomLog() {
	testFileFullPath := filepath.Join(suite.testDirName, "someFile.log")
	msLocation := proto.LogFileLocation{
		FileLocation: testFileFullPath,
		LogType:      proto.LogFileLocation_CUSTOM}
	_, e := FindLogFile(&msLocation)
	suite.Error(e, "File should not be found!")
}

func (suite *IntegrationFileOperationsTests) TestShouldFindMethodServerLogAsDefault() {
	msLocation := proto.LogFileLocation{FileLocation: suite.testDirName, LogType: 5}
	s, e := FindLogFile(&msLocation)
	if suite.NoError(e, "Default option should be found!") {
		suite.Equal(suite.correctFileNames[proto.LogFileLocation_METHOD_SERVER], s, "should be equal")
	}
}

func (suite *IntegrationFileOperationsTests) TestShouldFailWithIncorrectDirectory() {
	msLocation := proto.LogFileLocation{FileLocation: "incorrect", LogType: proto.LogFileLocation_METHOD_SERVER}
	s, e := FindLogFile(&msLocation)
	suite.Assertions.Error(e, "%s Should be an incorrect directory name", msLocation.FileLocation)
	suite.Assertions.Empty(s, "%s Should be empty", s)
}

func (suite *IntegrationFileOperationsTests) TestShouldFailWithNoFilesFound() {
	infos, _ := ioutil.ReadDir(suite.testDirName)
	for _, info := range infos {
		_ = os.Remove(suite.testDirName + "/" + info.Name())
	}

	msLocation := proto.LogFileLocation{FileLocation: suite.testDirName, LogType: proto.LogFileLocation_METHOD_SERVER}
	s, e := FindLogFile(&msLocation)
	suite.Assertions.EqualError(e, "Log file not found for type: METHOD_SERVER!",
		`Should not found any files`, msLocation.FileLocation)
	suite.Assertions.Empty(s, "%s Should be empty", s)
}

func removeFiles(s *IntegrationFileOperationsTests) {
	e := os.RemoveAll(s.testDirName)
	if e != nil {
		logger.Fatal(e)
	}
}

func createFiles(s *IntegrationFileOperationsTests) {
	s.testDirName, _ = ioutil.TempDir("", "temp_")
	createFile(filepath.Join(s.testDirName, s.correctFileNames[proto.LogFileLocation_METHOD_SERVER]))
	createFile(filepath.Join(s.testDirName, s.correctFileNames[proto.LogFileLocation_BACKGROUND_METHOD_SERVER]))
	createFile(filepath.Join(s.testDirName, s.correctFileNames[proto.LogFileLocation_CUSTOM]))
}

func createFile(name string) {
	f, e := os.OpenFile(name, os.O_CREATE, 777)
	if e != nil {
		logger.Fatal(e)
	}
	defer f.Close()
}
