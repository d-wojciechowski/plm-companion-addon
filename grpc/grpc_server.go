package grpc

import (
	"context"
	proto "dominikw.pl/wnc_plugin/proto"
	"dominikw.pl/wnc_plugin/util"
	"errors"
	"github.com/google/logger"
	"github.com/hpcloud/tail"
	"io/ioutil"
	"os"
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

func (s *Server) Navigate(ctx context.Context, protoPath *proto.Path) (*proto.FileResponse, error) {

	path := protoPath.Name
	if protoPath.Name == "" {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		path = filepath.Dir(ex)
	}

	fInfos, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var files []*proto.FileMeta
	for _, info := range fInfos {
		files = append(files, &proto.FileMeta{
			Name:        info.Name(),
			IsDirectory: info.IsDir(),
			Path:        path + string(os.PathSeparator) + info.Name(),
		})
	}

	return &proto.FileResponse{Metas: files}, err
}
