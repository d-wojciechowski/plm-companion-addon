package grpc

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestGRPCServer(t *testing.T) {
	suite.Run(t, new(gRpcLogsTestSuite))
	suite.Run(t, new(gRpcCommandTestSuite))
}
