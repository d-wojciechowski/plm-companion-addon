package main

import (
	grpcServer "dominikw.pl/wnc_plugin/grpc"
	proto "dominikw.pl/wnc_plugin/proto"
	"flag"
	"github.com/google/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"time"
)

var verbose = flag.Bool("v", false, "print info level logs to stdout")
var noWncMode = flag.Bool("noWnc", false, "turn on no wnc mode")

func main() {
	flag.Parse()
	defer setUpLogger().Close()

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		logger.Fatal(err)
	}

	server := grpc.NewServer()
	serviceServer := grpcServer.NewServer(*noWncMode)
	proto.RegisterCommandServiceServer(server, serviceServer)
	proto.RegisterLogViewerServiceServer(server, serviceServer)
	proto.RegisterFileServiceServer(server, serviceServer)
	reflection.Register(server)

	logger.Infof("server starting with parameters: -v: %t, -noWnc: %t", *verbose, *noWncMode)
	logger.Infof("server info: %v", listener.Addr())

	if e := server.Serve(listener); e != nil {
		logger.Fatal(e)
	}
}

func setUpLogger() *logger.Logger {
	_ = os.Mkdir("logs", os.ModeDir)
	_ = os.Chmod("logs", os.ModePerm)
	filename := "logs/" + time.Now().Format("2006_01_02-15_04") + ".log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error(err)
	}
	fileLogger := logger.Init(filename, *verbose, true, f)
	logger.SetFlags(log.LstdFlags)

	return fileLogger
}
