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

	tp := rsocket.TCPServer().SetAddr(srv.addr).Build()

	err := rsocket.Receive().
		Fragment(1024).
		Resume().
		Acceptor(func(ctx context.Context, setup payload.SetupPayload, socket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			logger.Info("Acceptor initialization started")
			defer logger.Info("Acceptor initialization ended")
			return rsocket.NewAbstractSocket(
				srv.requestResponseHandler(),
				srv.requestChannelHandler(),
				srv.requestStreamHandler(),
			), nil
		}).
		Transport(tp).
		Serve(ctx)
	util.PanicOnError(err)
}

func (srv *Server) requestResponseHandler() rsocket.OptAbstractSocket {
	logger.Info("RequestResponse initialization start")
	defer logger.Info("RequestResponse initialization start")
	return rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
		metadata, _ := msg.MetadataUTF8()
		if strings.Contains(metadata, constants_server.FileServiceIdentifier) &&
			strings.Contains(metadata, constants_server.NavigateIdentifier) {
			logger.Infof("Service %s with method %s execution start", constants_server.FileServiceIdentifier,
				constants_server.NavigateIdentifier)
			defer logger.Infof("Service %s with method %s execution ended", constants_server.FileServiceIdentifier,
				constants_server.NavigateIdentifier)
			return srv.Navigate(msg)
		}
		return srv.Execute(msg)
	})
}

func (srv *Server) requestChannelHandler() rsocket.OptAbstractSocket {
	logger.Infof("RequestChannel initialization start")
	defer logger.Infof("RequestChannel initialization ended")
	return rsocket.RequestChannel(func(requests flux.Flux) (responses flux.Flux) {
		return nil
	})
}

func (srv *Server) requestStreamHandler() rsocket.OptAbstractSocket {
	logger.Infof("requestStreamHandler initialization start")
	defer logger.Infof("requestStreamHandler initialization ended")
	return rsocket.RequestStream(func(msg payload.Payload) flux.Flux {
		metadata, _ := msg.MetadataUTF8()
		if strings.Contains(metadata, constants_server.LogServiceIdentifier) {
			logger.Infof("Service %s with method %s execution start", constants_server.LogServiceIdentifier,
				"GetLogs")
			defer logger.Infof("Service %s with method %s execution ended", constants_server.LogServiceIdentifier,
				"GetLogs")
			return srv.GetLogs(msg)
		}
		logger.Infof("Service %s with method %s execution start", constants_server.LogServiceIdentifier,
			"ExecuteStreaming")
		defer logger.Infof("Service %s with method %s execution ended", constants_server.LogServiceIdentifier,
			"ExecuteStreaming")
		return srv.ExecuteStreaming(msg)
	})
}

func toPayload(msg proto.Message, metadata []byte) payload.Payload {
	resp, _ := proto.Marshal(msg)
	return payload.New(resp, metadata)
}
