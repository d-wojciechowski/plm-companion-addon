package grpc

import (
	"bou.ke/monkey"
	"context"
	proto "dominikw.pl/wnc_plugin/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"os/exec"
	"reflect"
	"time"
)

type gRpcCommandTestSuite struct {
	suite.Suite
	server     *Server
	controller *gomock.Controller
}

// TODO create test cases for command-args pairs
func (s *gRpcCommandTestSuite) SetupSuite() {
	s.server = NewServer(false)
	s.controller = gomock.NewController(s.Suite.T())
	monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
		return &exec.Cmd{}
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(&exec.Cmd{}), "Run", func(c *exec.Cmd) error {
		return nil
	})
}

func (s *gRpcCommandTestSuite) TearDownSuite() {
	s.controller.Finish()
}

func (s *gRpcCommandTestSuite) TestServer_Execute() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.Run("testCorrectCall", testCorrectCall(s, ctx))
	s.Run("testCorrectCallWithNoArgs", testCorrectCallWithNoArgs(s, ctx))
	s.server.NoWncMode = true
	s.Run("testNoWncMode", testNoWncMode(s, ctx))

}

func testCorrectCallWithNoArgs(s *gRpcCommandTestSuite, ctx context.Context) func() {
	return func() {
		_, e := s.server.Execute(ctx, &proto.Command{
			Command: "windchill",
			Args:    "",
		})
		s.NoError(e)
	}
}

func testNoWncMode(s *gRpcCommandTestSuite, ctx context.Context) func() {
	return func() {
		response, e := s.server.Execute(ctx, &proto.Command{
			Command: "windchill",
			Args:    "stop",
		})

		if s.NoError(e) {
			s.Equal("NO WNC MODE", response.Message)
			s.Equal(int32(200), response.Status)
		}
	}
}

func testCorrectCall(s *gRpcCommandTestSuite, ctx context.Context) func() {
	return func() {
		_, e := s.server.Execute(ctx, &proto.Command{
			Command: "windchill",
			Args:    "stop",
		})

		s.NoError(e)
	}
}
