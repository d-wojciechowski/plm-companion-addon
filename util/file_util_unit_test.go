package util

import (
	proto "dominikw.pl/wnc_plugin/proto"
	"fmt"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

type UnitFileOperationsTests struct {
	suite.Suite
	correctFilenameForMS  string
	correctFilenameForBMS string
}

func TestCheckFileName(t *testing.T) {
	suite.Run(t, new(UnitFileOperationsTests))
}

func (s *UnitFileOperationsTests) BeforeTest(_, _ string) {
	now := time.Now()
	s.correctFilenameForMS = fmt.Sprintf("MethodServer1-%s%s%s319-30967-log4j.log",
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
	s.correctFilenameForBMS = fmt.Sprintf("BackgroundMethodServer1-%s%s%s319-30967-log4j.log",
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
}

func (s *UnitFileOperationsTests) TestCheckFileNameShouldWork() {
	s.True(checkFileName(s.correctFilenameForMS, proto.LogFileLocation_METHOD_SERVER), "%s Should be correct!", s.correctFilenameForMS)
}

func (s *UnitFileOperationsTests) TestCheckFileNameShouldFail() {
	name := "random test string"
	s.False(checkFileName(name, proto.LogFileLocation_METHOD_SERVER), "%s Should be incorrect!", name)
}

func (s *UnitFileOperationsTests) TestGetPathCustom() {
	path := GetPath("someDir", "test.log", &proto.LogFileLocation{LogType: proto.LogFileLocation_CUSTOM})
	s.Equal("test.log", path)
}

func (s *UnitFileOperationsTests) TestGetPathPredefined() {
	path := GetPath("someDir", "test.log", &proto.LogFileLocation{LogType: proto.LogFileLocation_METHOD_SERVER})
	s.Equal(filepath.Join("someDir", "test.log"), path)
}
