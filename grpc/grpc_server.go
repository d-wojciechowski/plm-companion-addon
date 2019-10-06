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
	"sort"
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
	logger.Infof("LogViewer request received, looking for log files in: %s", logFileDirectory)
	infos, e := ioutil.ReadDir(logFileDirectory)

	if e != nil {
		_ = outputStream.Send(&proto.LogLine{Message: "Directory not found!"})
		return e
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].ModTime().After(infos[j].ModTime())
	})

	logFileName := ""
	for _, info := range infos {
		if util.CheckFileName(info.Name()) {
			logFileName = info.Name()
			logger.Infof("Log file chosen: %s", logFileName)
			break
		}
		if util.CheckFileNameOmittingDate(info.Name()) {
			logFileName = info.Name()
			logger.Infof("Log file chosen: %s", logFileName)
			break
		}
	}

	if logFileName == "" {
		logger.Error("Log file not found!")
		return errors.New("log file not found")
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
