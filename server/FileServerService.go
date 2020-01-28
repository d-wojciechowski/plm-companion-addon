package server

import (
	"context"
	"dominikw.pl/wnc_plugin/proto/files"
	"dominikw.pl/wnc_plugin/util"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/google/logger"
	"github.com/hpcloud/tail"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func (s *Server) GetLogs(msg payload.Payload) flux.Flux {
	logFile := &files.LogFileLocation{}
	_ = proto.Unmarshal(msg.Data(), logFile)

	logFileDirectory := logFile.FileLocation
	logFileName, e := util.FindLogFile(logFile)
	if e != nil {
		flux.Error(e)
	}
	tailFile, e := tail.TailFile(util.GetPath(logFileDirectory, logFileName, logFile), s.tailConfig)
	if e != nil {
		logger.Error(e)
		flux.Error(e)
	}

	fluxResponse := flux.Create(func(ctx context.Context, s flux.Sink) {
		lines := tailFile.Lines
		for line := range lines {
			if e != nil {
				logger.Error(e)
				s.Error(e)
			}
			marshal, _ := proto.Marshal(&files.LogLine{Message: line.Text})
			s.Next(payload.New(marshal, nil))
		}
		s.Complete()
	}).DoFinally(func(s rx.SignalType) {
		_ = tailFile.Stop()
		tailFile.Cleanup()
	})

	return fluxResponse
}

func (srv *Server) Navigate(msg payload.Payload) mono.Mono {
	protoPath := &files.Path{}
	_ = proto.Unmarshal(msg.Data(), protoPath)

	paths := getPaths(protoPath)
	if len(paths) == 0 {
		mono.Error(errors.New("no path exception"))
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
	return mono.Just(toPayload(getFullResult(root, protoPath.FullExpand), make([]byte, 1)))
}

func fillAncestor(ancestor *files.FileMeta, currentPath string, nextElement string) (intermediateAncestor *files.FileMeta) {
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

func getFullResult(root *files.FileMeta, fullExpand bool) *files.FileResponse {
	result := []*files.FileMeta{root}
	if fullExpand {
		result = addOtherDrives(result)
	}
	return &files.FileResponse{
		FileTree:  result,
		Separator: string(os.PathSeparator),
		Os:        runtime.GOOS,
	}
}

func getPaths(protoPath *files.Path) []string {
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

func addOtherDrives(result []*files.FileMeta) []*files.FileMeta {
	elements := funk.Filter(util.GetWindowsDrives(), func(s string) bool {
		return string(result[0].GetName()[0]) != s
	})
	for _, drive := range elements.([]string) {
		result = append(result, buildFileMeta(drive+":", true))
	}
	return result
}

func buildFileMeta(name string, IsDir bool) *files.FileMeta {
	return &files.FileMeta{
		Name:        name,
		IsDirectory: IsDir,
		ChildFiles:  []*files.FileMeta{},
	}
}
