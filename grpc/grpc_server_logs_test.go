package grpc

import (
	"bou.ke/monkey"
	"context"
	projectMock "dominikw.pl/wnc_plugin/mock"
	proto "dominikw.pl/wnc_plugin/proto"
	"dominikw.pl/wnc_plugin/util"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/logger"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"reflect"
	"time"
)

type gRpcLogsTestSuite struct {
	suite.Suite
	server      *Server
	controller  *gomock.Controller
	testDirName string
}

func (s *gRpcLogsTestSuite) SetupSuite() {
	var e error
	s.controller = gomock.NewController(s.Suite.T())
	s.server = NewServer(false)
	// Disable log follow to avoid infinite stream
	s.server.tailConfig.ReOpen = false
	s.server.tailConfig.Follow = false
	s.testDirName, e = ioutil.TempDir("", "go_test_temp_")
	if e != nil {
		logger.Fatal(e)
	}
	f, e := ioutil.TempFile(s.testDirName, "MethodServer-1910071229-*-log4j.log")
	defer f.Close()
	if e != nil {
		logger.Fatal(e)
	}
	for i := 0; i < 20; i++ {
		_, e := f.WriteString(fmt.Sprintf("Some Content #%d \n", i))
		if e != nil {
			logger.Fatal(e)
		}
	}
}

func (s *gRpcLogsTestSuite) TearDownSuite() {
	s.controller.Finish()
}

func (s *gRpcLogsTestSuite) TestServer_GetLogs() {
	s.Run("testCorrectLogRead", s.testCorrectLogRead)
	s.Run("testFileNotFound", s.testFileNotFound)
	s.Run("testErrorInTail", s.testErrorInTail)
	s.Run("testIncorrectLogSend", s.testIncorrectLogSend)
	s.Run("testClientInterrupted", s.testClientInterrupted)
}

func (s *gRpcLogsTestSuite) testIncorrectLogSend() {
	serviceGetLogsServer := projectMock.NewMockLogViewerService_GetLogsServer(s.controller)
	serviceGetLogsServer.EXPECT().
		Context().
		Return(context.Background())

	serviceGetLogsServer.EXPECT().
		Send(gomock.AssignableToTypeOf(&proto.LogLine{})).
		Return(errors.New("some IO error"))

	location := &proto.LogFileLocation{
		FileLocation: s.testDirName,
		LogType:      proto.LogFileLocation_METHOD_SERVER,
	}

	e := s.server.GetLogs(location, serviceGetLogsServer)
	s.Error(e, "some IO error")
}

func (s *gRpcLogsTestSuite) testCorrectLogRead() {
	serviceGetLogsServer := projectMock.NewMockLogViewerService_GetLogsServer(s.controller)

	serviceGetLogsServer.EXPECT().
		Context().
		Return(context.TODO())

	serviceGetLogsServer.EXPECT().
		Send(gomock.AssignableToTypeOf(&proto.LogLine{})).
		AnyTimes().
		Return(nil)

	location := &proto.LogFileLocation{
		FileLocation: s.testDirName,
		LogType:      proto.LogFileLocation_METHOD_SERVER,
	}

	e := s.server.GetLogs(location, serviceGetLogsServer)
	s.NoError(e)
}

func (s *gRpcLogsTestSuite) testFileNotFound() {
	serviceGetLogsServer := projectMock.NewMockLogViewerService_GetLogsServer(s.controller)

	serviceGetLogsServer.EXPECT().
		Send(gomock.AssignableToTypeOf(&proto.LogLine{})).
		AnyTimes().
		Return(nil)

	location := &proto.LogFileLocation{
		FileLocation: "asd",
		LogType:      proto.LogFileLocation_METHOD_SERVER,
	}

	e := s.server.GetLogs(location, serviceGetLogsServer)
	if s.Error(e, "File should not be found") {
		s.IsType(&os.PathError{}, e, "Wrong error type")
	}
}

func (s *gRpcLogsTestSuite) testErrorInTail() {
	serviceGetLogsServer := projectMock.NewMockLogViewerService_GetLogsServer(s.controller)

	serviceGetLogsServer.EXPECT().
		Send(gomock.AssignableToTypeOf(&proto.LogLine{})).
		AnyTimes().
		Return(nil)

	location := &proto.LogFileLocation{
		FileLocation: s.testDirName,
		LogType:      proto.LogFileLocation_METHOD_SERVER,
	}

	monkey.Patch(util.FindLogFile, func(*proto.LogFileLocation) (string, error) {
		return "wrongFilePath", nil
	})

	e := s.server.GetLogs(location, serviceGetLogsServer)
	s.Error(e)
	monkey.UnpatchAll()
}

func (s *gRpcLogsTestSuite) testClientInterrupted() {
	serviceGetLogsServer := projectMock.NewMockLogViewerService_GetLogsServer(s.controller)
	s.server.tailConfig.Follow = true
	timeout, cancel := context.WithCancel(context.Background())
	serviceGetLogsServer.EXPECT().
		Context().
		Return(timeout)

	go func() {
		select {
		case <-time.After(5 * time.Second):
			cancel()
			s.Fail("gorutine stuck")
		}
	}()

	serviceGetLogsServer.EXPECT().
		Send(gomock.AssignableToTypeOf(&proto.LogLine{})).
		MinTimes(1).
		MaxTimes(4).
		Return(nil)

	location := &proto.LogFileLocation{
		FileLocation: s.testDirName,
		LogType:      proto.LogFileLocation_METHOD_SERVER,
	}

	monkey.PatchInstanceMethod(reflect.TypeOf(serviceGetLogsServer), "Send",
		func(_ *projectMock.MockLogViewerService_GetLogsServer, _ *proto.LogLine) error {
			cancel()
			return nil
		})

	e := s.server.GetLogs(location, serviceGetLogsServer)
	s.NoError(e)
	monkey.UnpatchAll()
}
