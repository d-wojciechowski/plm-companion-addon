package main

import (
	"context"
	"dominikw.pl/wnc_plugin/proto"
	"github.com/hpcloud/tail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os/exec"
)

type server struct{}

func main() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterCommandServiceServer(srv, &server{})
	proto.RegisterLogViewerServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) Execute(ctx context.Context, command *proto.Command) (*proto.Response, error) {
	var cmd *exec.Cmd
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
	tailFile, send := tail.TailFile("test.log", config)

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
