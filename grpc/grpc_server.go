package grpc

import (
	"context"
	proto "dominikw.pl/wnc_plugin/proto"
	"dominikw.pl/wnc_plugin/util"
	"github.com/google/logger"
	"github.com/hpcloud/tail"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

	paths := strings.Split(path, string(os.PathSeparator))

	cuurrentPath := ""
	var root *proto.FileMeta
	var ancestor *proto.FileMeta

	for i := 0; i < len(paths); i++ {
		if root == nil {
			root = &proto.FileMeta{
				Name:        paths[i],
				IsDirectory: true,
				ChildFiles:  []*proto.FileMeta{},
			}
			ancestor = root
		}
		if cuurrentPath == "" {
			cuurrentPath += paths[i] + string(os.PathSeparator)
		} else {
			cuurrentPath += string(os.PathSeparator) + paths[i]
		}
		var intermediateAncestor *proto.FileMeta
		fInfos, _ := ioutil.ReadDir(cuurrentPath)
		for _, info := range fInfos {
			currentFM := &proto.FileMeta{
				Name:        info.Name(),
				IsDirectory: info.IsDir(),
				ChildFiles:  []*proto.FileMeta{},
			}
			ancestor.ChildFiles = append(ancestor.ChildFiles, currentFM)
			if len(paths) > i+1 && info.Name() == paths[i+1] {
				intermediateAncestor = currentFM
			}
		}
		ancestor = intermediateAncestor
	}

	return &proto.FileResponse{FileTree: []*proto.FileMeta{root}}, nil
}
