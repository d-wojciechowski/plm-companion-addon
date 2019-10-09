package grpc

import (
	"context"
	proto "dominikw.pl/wnc_plugin/proto"
	"dominikw.pl/wnc_plugin/util"
	"github.com/google/logger"
	"github.com/hpcloud/tail"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Server struct {
	NoWncMode  bool
	tailConfig tail.Config
}

func NewServer(noWnc bool) *Server {
	return &Server{
		NoWncMode: noWnc,
		tailConfig: tail.Config{
			ReOpen:    true,
			MustExist: true,
			Follow:    true,
			Poll:      runtime.GOOS == "windows",
		},
	}
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

	logFileDirectory := logFile.FileLocation
	logFileName, i := util.FindLogFile(logFile)
	if i != nil {
		return i
	}

	tailFile, send := tail.TailFile(filepath.Join(logFileDirectory, logFileName), s.tailConfig)

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
