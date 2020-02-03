package server

import (
	"dominikw.pl/wnc_plugin/proto/commands"
	constants_messages "dominikw.pl/wnc_plugin/server/constants/messages"
	"dominikw.pl/wnc_plugin/util"
	"github.com/golang/protobuf/proto"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
	"os/exec"
)

func (srv *Server) Execute(msg payload.Payload) mono.Mono {
	command := &commands.Command{}
	_ = proto.Unmarshal(msg.Data(), command)

	var cmd *exec.Cmd
	if srv.NoWncMode {
		return mono.Just(toPayload(&commands.Response{Message: constants_messages.NoWncMode, Status: 200}, make([]byte, 1)))
	}

	if !util.IsEmpty(command.GetArgs()) {
		cmd = exec.Command(command.GetCommand(), command.GetArgs())
	} else {
		cmd = exec.Command(command.GetCommand())
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return mono.Error(err)
	}

	return mono.Just(toPayload(&commands.Response{Message: string(out), Status: 200}, make([]byte, 1)))
}
