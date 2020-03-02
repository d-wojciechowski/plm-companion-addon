package main

import (
	"dominikw.pl/wnc_plugin/server"
	"flag"
	"github.com/google/logger"
	"log"
	"os"
	"strconv"
	"time"
)

var verbose = flag.Bool("v", false, "print info level logs to stdout")
var noWncMode = flag.Bool("noWnc", false, "turn on no wnc mode")
var portNumber = flag.Int("port", 4040, "port number on which server will be listening.")

func main() {
	flag.Parse()
	defer setUpLogger().Close()

	logger.Infof("server starting with parameters: -v: %t, -noWnc: %t, -port: %d", *verbose, *noWncMode, *portNumber)

	rsocketServer := server.NewServer(*noWncMode, "tcp://0.0.0.0:"+strconv.Itoa(*portNumber))
	go func() { rsocketServer.Start() }()

	select {}
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
