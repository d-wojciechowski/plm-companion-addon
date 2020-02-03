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
	"time"
)

func (srv *Server) Execute(msg payload.Payload) mono.Mono {
	command := &commands.Command{}
	_ = proto.Unmarshal(msg.Data(), command)

	var cmd *exec.Cmd
	if srv.NoWncMode {
		return mono.Just(toPayload(&commands.Response{Message: "NO WNC MODE", Status: commands.Status_FINISHED}, make([]byte, 1)))
	}

	if command.GetArgs() != "" {
		cmd = exec.Command(command.GetCommand(), command.GetArgs())
	} else {
		cmd = exec.Command(command.GetCommand())
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return mono.Error(err)
	}

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
		var cmd *exec.Cmd
		if command.GetArgs() != "" {
			cmd = exec.Command(command.GetCommand(), command.GetArgs())
		} else {
			cmd = exec.Command(command.GetCommand())
		}

		stdout, err := cmd.StdoutPipe()
		errorPiper, err := cmd.StderrPipe()

		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				m := scanner.Text()
				s.Next(toPayload(&commands.Response{Message: m, Status: commands.Status_RUNNING}, make([]byte, 1)))
			}
			s.Complete()
		}()

		go func() {
			scanner := bufio.NewScanner(errorPiper)
			for scanner.Scan() {
				s.Error(errors.New(scanner.Text()))
			}
		}()

		if err != nil {
			s.Error(err)
		}

	})
}
