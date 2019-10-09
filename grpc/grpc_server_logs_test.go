package grpc

import (
	projectMock "dominikw.pl/wnc_plugin/mock"
	proto "dominikw.pl/wnc_plugin/proto"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/logger"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
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
	s.testDirName, e = ioutil.TempDir("", "temp_")
	if e != nil {
		logger.Fatal(e)
	}
	f, e := ioutil.TempFile(s.testDirName, "MethodServer-1910071229-3412-log4j.log")
	if e != nil {
		logger.Fatal(e)
	}
	_, _ = f.WriteString("some content")
	defer f.Close()
}

func (s *gRpcLogsTestSuite) TearDownSuite() {
	s.controller.Finish()
	all := os.RemoveAll(s.testDirName)
	if all != nil {
		logger.Fatal(all)
	}
}

func (s *gRpcLogsTestSuite) TestServer_GetLogs() {

	s.Run("testCorrectLogRead", s.testCorrectLogRead)
	s.Run("testIncorrectLogRead", s.testIncorrectLogRead)

}

func (s *gRpcLogsTestSuite) testIncorrectLogRead() {
	serviceGetLogsServer := projectMock.NewMockLogViewerService_GetLogsServer(s.controller)

	serviceGetLogsServer.EXPECT().
		Send(gomock.AssignableToTypeOf(&proto.LogLine{})).
		Return(errors.New("some IO error"))

	location := &proto.LogFileLocation{
		FileLocation: s.testDirName,
		LogType:      proto.LogFileLocation_METHOD_SERVER,
	}

	e := s.server.GetLogs(location, serviceGetLogsServer)
	s.EqualError(e, "some IO error")
}

func (s *gRpcLogsTestSuite) testCorrectLogRead() {
	serviceGetLogsServer := projectMock.NewMockLogViewerService_GetLogsServer(s.controller)

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
