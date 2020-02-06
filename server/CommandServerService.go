package server

import (
	"bufio"
	"context"
	"dominikw.pl/wnc_plugin/proto/commands"
	"github.com/golang/protobuf/proto"
	"github.com/google/logger"
	"github.com/pkg/errors"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

func (srv *Server) Execute(msg payload.Payload) mono.Mono {
	command := &commands.Command{}
	_ = proto.Unmarshal(msg.Data(), command)

	if srv.NoWncMode {
		return mono.Just(toPayload(&commands.Response{Message: "NO WNC MODE", Status: commands.Status_FINISHED}, make([]byte, 1)))
	}

	if command.GetArgs() != "" {
		command.Args = ""
	}
	cmd := execCommand(command)
	out, _ := cmd.CombinedOutput()
	return mono.Just(toPayload(&commands.Response{Message: string(out), Status: commands.Status_FINISHED}, make([]byte, 1)))
}

func (srv *Server) ExecuteStreaming(msg payload.Payload) flux.Flux {
	command := &commands.Command{}
	_ = proto.Unmarshal(msg.Data(), command)

	if srv.NoWncMode {
		return flux.Create(func(ctx context.Context, s flux.Sink) {
			for i := 0; i < 4; i++ {
				s.Next(toPayload(&commands.Response{Message: "NO WNC MODE", Status: commands.Status_FINISHED}, make([]byte, 1)))
				time.Sleep(2 * time.Second)
			}
			s.Complete()
		})
	}

	var cmd *exec.Cmd
	return flux.Create(func(ctx context.Context, s flux.Sink) {
		if command.GetArgs() != "" {
			command.Args = ""
		}
		cmd = execCommand(command)

		stdout, err := cmd.StdoutPipe()
		errorPiper, err := cmd.StderrPipe()

		go func() {
			scanner := bufio.NewScanner(stdout)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				m := scanner.Text()
				s.Next(toPayload(&commands.Response{Message: m, Status: commands.Status_RUNNING}, make([]byte, 1)))
			}
			s.Complete()
		}()

		go func() {
			scanner := bufio.NewScanner(errorPiper)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				s.Error(errors.New(scanner.Text()))
			}
		}()

		_ = cmd.Start()
		if err != nil {
			s.Error(err)
		}

	}).DoFinally(func(s rx.SignalType) {
		kill(cmd)
	})
}

func kill(cmd *exec.Cmd) {
	var err error
	if runtime.GOOS == "windows" {
		err = killOnWindows(cmd)
	} else {
		err = killOnLinux(cmd)
	}
	logger.Error(err)
}

func killOnLinux(cmd *exec.Cmd) error {
	kill := exec.Command("pkill", "-P", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}

func killOnWindows(cmd *exec.Cmd) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}

func execCommand(cmd *commands.Command) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/U", "/c", cmd.GetCommand())
	} else {
		return exec.Command("sh", "-c", cmd.GetCommand())
	}
}
