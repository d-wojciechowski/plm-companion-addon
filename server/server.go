package server

import (
	"context"
	"github.com/d-wojciechowski/plm-companion-addon/server/constants/other"
	"github.com/d-wojciechowski/plm-companion-addon/server/constants/server"
	"github.com/d-wojciechowski/plm-companion-addon/util"
	"github.com/golang/protobuf/proto"
	"github.com/google/logger"
	"github.com/hpcloud/tail"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"runtime"
	"strings"
)

type Server struct {
	addr       string
	devMode    bool
	tailConfig tail.Config
}

func NewServer(devMode bool, addr string) *Server {
	logger.Infof("Attempt to create server on addr %s", addr)
	defer logger.Infof("Server instance started on addr %s", addr)
	return &Server{
		addr:    addr,
		devMode: devMode,
		tailConfig: tail.Config{
			ReOpen:    true,
			MustExist: true,
			Follow:    true,
			Poll:      runtime.GOOS == constants_other.WindowsOSName,
		},
	}
}

func (srv *Server) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := rsocket.Receive().
		Fragment(1024).
		Resume().
		Acceptor(func(setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			return rsocket.NewAbstractSocket(
				srv.requestResponseHandler(),
				srv.requestChannelHandler(),
				srv.requestStreamHandler(),
			), nil
		}).
		Transport(srv.addr).
		Serve(ctx)
	util.PanicOnError(err)
}

func (srv *Server) requestResponseHandler() rsocket.OptAbstractSocket {
	return rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
		metadata, _ := msg.MetadataUTF8()
		if strings.Contains(metadata, constants_server.FileServiceIdentifier) &&
			strings.Contains(metadata, constants_server.NavigateIdentifier) {
			return srv.Navigate(msg)
		}
		return srv.Execute(msg)
	})
}

func (srv *Server) requestChannelHandler() rsocket.OptAbstractSocket {
	return rsocket.RequestChannel(func(msgs rx.Publisher) flux.Flux {
		return nil
	})
}

func (srv *Server) requestStreamHandler() rsocket.OptAbstractSocket {
	return rsocket.RequestStream(func(msg payload.Payload) flux.Flux {
		metadata, _ := msg.MetadataUTF8()
		if strings.Contains(metadata, "LogViewerService") {
			return srv.GetLogs(msg)
		}
		return srv.ExecuteStreaming(msg)
	})
}

func toPayload(msg proto.Message, metadata []byte) payload.Payload {
	resp, _ := proto.Marshal(msg)
	return payload.New(resp, metadata)
}
