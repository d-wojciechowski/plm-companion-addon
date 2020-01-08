package grpc

import (
	"context"
	proto "dominikw.pl/wnc_plugin/proto"
	"dominikw.pl/wnc_plugin/util"
	"errors"
	"github.com/google/logger"
	"github.com/hpcloud/tail"
	"github.com/thoas/go-funk"
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

func (s *Server) GetLogs(logFile *proto.LogFileLocation, outputStream proto.LogViewerService_GetLogsServer) (e error) {

	logFileDirectory := logFile.FileLocation
	logFileName, e := util.FindLogFile(logFile)
	if e != nil {
		return
	}
	tailFile, e := tail.TailFile(util.GetPath(logFileDirectory, logFileName, logFile), s.tailConfig)

	if e != nil {
		logger.Error(e)
		return
	}
	lines := tailFile.Lines
	go handleShutdownByClient(outputStream.Context(), tailFile)
	for line := range lines {
		e = outputStream.Send(&proto.LogLine{Message: line.Text})
		if e != nil {
			logger.Error(e)
			break
		}
	}

	return
}

func handleShutdownByClient(ctx context.Context, tailFile *tail.Tail) {
	<-ctx.Done()
	_ = tailFile.Stop()
	tailFile.Cleanup()
	logger.Warning("Interrupted by client")
}

/*Navigate retuns File structure starting from path, if proto.Path states that path should be fully expanded, it traverse starting from root */
func (s *Server) Navigate(ctx context.Context, protoPath *proto.Path) (*proto.FileResponse, error) {
	paths := getPaths(protoPath)
	if len(paths) == 0 {
		return nil, errors.New("no path exception")
	}
	currentPath := ""
	root := buildFileMeta(paths[0], true)
	ancestor := root
	for index, elem := range paths {
		currentPath = getCurrentPath(currentPath, elem)
		var nextElement string
		if len(paths) > index+1 {
			nextElement = paths[index+1]
		}
		ancestor = fillAncestor(ancestor, currentPath, nextElement)
	}
	return getFullResult(root, protoPath.FullExpand), nil
}

func fillAncestor(ancestor *proto.FileMeta, currentPath string, nextElement string) (intermediateAncestor *proto.FileMeta) {
	fInfos, e := ioutil.ReadDir(currentPath)
	if e != nil {
		logger.Error(e)
	}
	for _, info := range fInfos {
		currentFM := buildFileMeta(info.Name(), info.IsDir())
		ancestor.ChildFiles = append(ancestor.ChildFiles, currentFM)
		if nextElement != "" && info.Name() == nextElement {
			intermediateAncestor = currentFM
		}
	}
	return
}

func getCurrentPath(currentPath string, elem string) string {
	if currentPath == "" {
		return elem + string(os.PathSeparator)
	}
	return currentPath + string(os.PathSeparator) + elem
}

func getFullResult(root *proto.FileMeta, fullExpand bool) *proto.FileResponse {
	result := []*proto.FileMeta{root}
	if fullExpand {
		result = addOtherDrives(result)
	}
	return &proto.FileResponse{
		FileTree:  result,
		Separator: string(os.PathSeparator),
		Os:        runtime.GOOS,
	}
}

func getPaths(protoPath *proto.Path) []string {
	path := protoPath.Name
	if protoPath.Name == "" {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		path = filepath.Dir(ex)
		if runtime.GOOS == "windows" {
			path = strings.ToUpper(path[:1]) + path[1:]
		}
	}
	var paths []string
	if protoPath.FullExpand {
		paths = strings.Split(path, string(os.PathSeparator))
		if runtime.GOOS != "windows" {
			paths = append([]string{"/"}, paths...)
		}
	} else {
		paths = []string{path}
	}

	return funk.Filter(paths, func(s string) bool {
		return s != ""
	}).([]string)
}

func addOtherDrives(result []*proto.FileMeta) []*proto.FileMeta {
	elements := funk.Filter(util.GetWindowsDrives(), func(s string) bool {
		return string(result[0].GetName()[0]) != s
	})
	for _, drive := range elements.([]string) {
		result = append(result, buildFileMeta(drive+":", true))
	}
	return result
}

func buildFileMeta(name string, IsDir bool) *proto.FileMeta {
	return &proto.FileMeta{
		Name:        name,
		IsDirectory: IsDir,
		ChildFiles:  []*proto.FileMeta{},
	}
}
