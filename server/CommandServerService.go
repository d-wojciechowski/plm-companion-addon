package server

import (
	"bufio"
	"context"
	"github.com/d-wojciechowski/plm-companion-addon/proto/commands"
	constants_messages "github.com/d-wojciechowski/plm-companion-addon/server/constants/messages"
	constants_other "github.com/d-wojciechowski/plm-companion-addon/server/constants/other"
	"github.com/golang/protobuf/proto"
	"github.com/google/logger"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func (srv *Server) Execute(msg payload.Payload) mono.Mono {
	command := &commands.Command{}
	_ = proto.Unmarshal(msg.Data(), command)

	if srv.devMode {
		return mono.Just(toPayload(&commands.Response{Message: constants_messages.NoWncMode, Status: commands.Status_FINISHED}, make([]byte, 1)))
	}

	cmd := execCommand(command)
	out, _ := cmd.CombinedOutput()
	return mono.Just(toPayload(&commands.Response{Message: string(out), Status: commands.Status_FINISHED}, make([]byte, 1)))
}

func (srv *Server) ExecuteStreaming(msg payload.Payload) flux.Flux {
	command := &commands.Command{}
	_ = proto.Unmarshal(msg.Data(), command)
	logger.Infof("Execution of command %s started", command.Command)

	if srv.devMode {
		logger.Infof("Dummy dev mode execution of command %s", command.Command)
		return flux.Create(func(ctx context.Context, s flux.Sink) {
			for i := 0; i < 4; i++ {
				response := &commands.Response{
					Message: constants_messages.NoWncMode,
					Status:  commands.Status_FINISHED}
				logger.Infof("Dummy message: %s with status %s", constants_messages.NoWncMode, commands.Status_FINISHED)
				s.Next(toPayload(response, make([]byte, 1)))
				time.Sleep(2 * time.Second)
			}
			logger.Infof("Execution of command %s ended", command.Command)
			s.Complete()
		})
	}

	var cmd *exec.Cmd
	return flux.Create(func(ctx context.Context, s flux.Sink) {
		cmd = execCommand(command)

		stdout, err := cmd.StdoutPipe()
		errorPiper, err := cmd.StderrPipe()

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			logger.Infof("Command %s standard out opened", command.Command)
			pipeReader(stdout, s, commands.Status_RUNNING)
			wg.Done()
			logger.Infof("Command %s standard out closed", command.Command)
		}()

		go func() {
			logger.Infof("Command %s error out opened", command.Command)
			pipeReader(errorPiper, s, commands.Status_FAILED)
			wg.Done()
			logger.Infof("Command %s error out closed", command.Command)
		}()

		go func() {
			logger.Infof("Command %s will run until both outs will finish", command.Command)
			wg.Wait()
			logger.Infof("Command %s stdout and err out closed", command.Command)
			s.Complete()
			logger.Infof("Command %s completed", command.Command)
		}()

		_ = cmd.Start()
		if err != nil {
			logger.Errorf("Command %s finished with error %s", command.Command, err.Error())
			s.Error(err)
		}

	}).DoFinally(func(s rx.SignalType) {
		kill(cmd)
	})
}

func pipeReader(pipe io.ReadCloser, s flux.Sink, status commands.Status) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		s.Next(toPayload(&commands.Response{Message: m, Status: status}, make([]byte, 1)))
	}
}

func kill(cmd *exec.Cmd) {
	logger.Info("Attempt to kill")
	if runtime.GOOS == constants_other.WindowsOSName {
		killOnWindows(cmd)
	} else {
		killOnLinux(cmd)
	}
	logger.Info("Kill success")
}

func killOnLinux(cmd *exec.Cmd) {
	kill := exec.Command("pkill", "-P", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	_ = kill.Run()
}

func killOnWindows(cmd *exec.Cmd) {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	_ = kill.Run()
}

func execCommand(cmd *commands.Command) *exec.Cmd {
	if runtime.GOOS == constants_other.WindowsOSName {
		return exec.Command("cmd", "/U", "/c", cmd.GetCommand())
	} else {
		return exec.Command("sh", "-c", cmd.GetCommand())
	}
}
