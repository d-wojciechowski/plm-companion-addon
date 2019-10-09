package grpc

import (
	projectMock "dominikw.pl/wnc_plugin/mock"
	proto "dominikw.pl/wnc_plugin/proto"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type gRpcLogsTestSuite struct {
	suite.Suite
	server     *Server
	controller *gomock.Controller
}

func (s *gRpcLogsTestSuite) SetupSuite() {
	s.controller = gomock.NewController(s.Suite.T())
	s.server = NewServer(false)
	// Disable log follow to avoid infinite stream
	s.server.tailConfig.ReOpen = false
	s.server.tailConfig.Follow = false
}

func (s *gRpcLogsTestSuite) TearDownSuite() {
	s.controller.Finish()
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
		FileLocation: "..",
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
		FileLocation: "..",
		LogType:      proto.LogFileLocation_METHOD_SERVER,
	}

	e := s.server.GetLogs(location, serviceGetLogsServer)
	s.NoError(e)
}
