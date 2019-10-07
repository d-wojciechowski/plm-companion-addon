package grpc

import (
	"context"
	proto "dominikw.pl/wnc_plugin/proto"
	"dominikw.pl/wnc_plugin/util"
	"github.com/google/logger"
	"github.com/hpcloud/tail"
	"os/exec"
	"path/filepath"
)

type Server struct {
	NoWncMode bool
}

func (s *Server) Execute(ctx context.Context, command *proto.Command) (*proto.Response, error) {
	var cmd *exec.Cmd

	if s.NoWncMode {
		return &proto.Response{Message: "NO WNC MODE", Status: 200}, nil
	}

	if command.GetArgs() != "" {
		cmd = exec.Command(command.GetCommand(), command.GetArgs())
	} else {
		cmd = exec.Command(command.GetCommand())
	}
	out, err := cmd.CombinedOutput()
	return &proto.Response{Message: string(out), Status: 200}, err
}

func (s *Server) GetLogs(logFile *proto.LogFileLocation, outputStream proto.LogViewerService_GetLogsServer) error {

	config := tail.Config{
		ReOpen:    true,
		MustExist: true,
		Follow:    true,
	}
	logFileDirectory := logFile.FileLocation
	logFileName, i := util.FindLogFile(logFile, outputStream)
	if i != nil {
		return i
	}

	tailFile, send := tail.TailFile(filepath.Join(logFileDirectory, logFileName), config)

	if send != nil {
		logger.Error(send)
		return send
	}

	lines := tailFile.Lines
	for line := range lines {
		send = outputStream.Send(&proto.LogLine{Message: line.Text})
		if send != nil {
			logger.Error(send)
			return send
		}
	}

	return nil
}
