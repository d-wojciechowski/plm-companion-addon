package server

import (
	"bufio"
	"context"
	"dominikw.pl/wnc_plugin/proto/commands"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"os/exec"
	"runtime"
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

	return flux.Create(func(ctx context.Context, s flux.Sink) {
		if command.GetArgs() != "" {
			command.Args = ""
		}
		cmd := execCommand(command)

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

	})
}

func execCommand(cmd *commands.Command) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/U", "/c", cmd.GetCommand(), cmd.GetArgs())
	} else {
		return exec.Command("sh", "-c", cmd.GetCommand(), cmd.GetArgs())
	}
}
