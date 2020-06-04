package main

import (
	"flag"
	"fmt"
	"github.com/d-wojciechowski/plm-companion-addon/server"
	"github.com/google/logger"
	"log"
	"os"
	"strconv"
	"time"
)

var verbose = flag.Bool("v", false, "Print info level logs to stdout")
var portNumber = flag.Int("port", 4040, "Port number on which server will be listening.")
var devMode = flag.Bool("devMode", false, "Turn on dev mode.")

func main() {
	flag.Parse()
	defer setUpLogger().Close()

	fmt.Printf("Server starting with parameters: -v: %t, -noWnc: %t, -port: %d\n", *verbose, *devMode, *portNumber)
	logger.Infof("Server starting with parameters: -v: %t, -noWnc: %t, -port: %d", *verbose, *devMode, *portNumber)

	go func() { server.NewServer(*devMode, "0.0.0.0:"+strconv.Itoa(*portNumber)).Start() }()

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
