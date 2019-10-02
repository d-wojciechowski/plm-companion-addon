package main

import (
	"context"
	"github.com/hpcloud/tail"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"dominikw.pl/wnc_plugin/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

var fileLogger *log.Logger

var logDisabled bool
var noWncMode bool

func main() {
	getConfig()
	defer setUpLogger().Close()

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		fileLogger.Panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterCommandServiceServer(srv, &server{})
	proto.RegisterLogViewerServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		fileLogger.Panic(e)
	}
}

func (s *server) Execute(ctx context.Context, command *proto.Command) (*proto.Response, error) {
	var cmd *exec.Cmd

	if !logDisabled {
		fileLogger.Println(command.Command + " " + command.Args)
	}
	if noWncMode {
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

func (s *server) GetLogs(logFile *proto.LogFileLocation, outputStream proto.LogViewerService_GetLogsServer) error {

	config := tail.Config{
		ReOpen:    true,
		MustExist: true,
		Follow:    true,
	}
	tailFile, send := tail.TailFile(logFile.FileLocation, config)

	if send != nil {
		return send
	}

	lines := tailFile.Lines
	for line := range lines {
		send = outputStream.Send(&proto.LogLine{Message: line.Text})
		if send != nil {
			return send
		}
	}

	return nil
}

func getConfig() {
	for _, s := range os.Args {
		if s == "-noLog" {
			logDisabled = true
		} else if s == "-noWnc" {
			noWncMode = true
		}
	}
}

func setUpLogger() *os.File {
	_ = os.Mkdir("logs", os.ModeDir)
	_ = os.Chmod("logs", os.ModePerm)
	filename := "logs/" + time.Now().Format("2006_01_02-15_04") + ".log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	fileLogger = log.New(f, "", log.LstdFlags)
	return f
}
