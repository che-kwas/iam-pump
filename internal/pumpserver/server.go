package pumpserver

import (
	"context"
	"iam-pump/internal/pumpserver/consumer"
	"iam-pump/internal/pumpserver/consumer/kafka"
	"iam-pump/internal/pumpserver/pump"
	"iam-pump/internal/pumpserver/store"
	"iam-pump/internal/pumpserver/store/mongo"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/che-kwas/iam-kit/server"
	"github.com/che-kwas/iam-kit/shutdown"
)

type pumpServer struct {
	*server.Server
	name   string
	ctx    context.Context
	cancel context.CancelFunc
	log    *logger.Logger

	err error
}

// NewServer builds a new pumpServer.
func NewServer(name string) *pumpServer {
	ctx, cancel := context.WithCancel(context.Background())

	s := &pumpServer{
		name:   name,
		ctx:    ctx,
		cancel: cancel,
		log:    logger.L(),
	}

	return s.initStore().initPump().newServer()
}

// Run runs the pumpServer.
func (s *pumpServer) Run() {
	if s.err != nil {
		s.log.Fatal(s.err)
	}

	defer s.cancel()
	defer s.log.Sync()

	if err := s.Server.Run(); err != nil {
		s.log.Fatal(err)
	}
}

func (s *pumpServer) initStore() *pumpServer {
	var cli store.Store
	if cli, s.err = mongo.NewMongoStore(); s.err != nil {
		return s
	}
	store.SetClient(cli)

	return s
}

func (s *pumpServer) initPump() *pumpServer {
	if s.err != nil {
		return s
	}

	var c consumer.Consumer
	msgHandler := consumer.MsgHandler(pump.TransferAuditRecord)
	c, s.err = kafka.NewConsumer(s.ctx, msgHandler)
	if s.err != nil {
		return s
	}
	consumer.SetConsumer(c)

	go pump.InitPump(s.ctx, c).Start()
	return s
}

func (s *pumpServer) newServer() *pumpServer {
	if s.err != nil {
		return s
	}

	s.Server, s.err = server.NewServer(
		s.name,
		server.WithShutdown(shutdown.ShutdownFunc(store.Client().Close)),
		server.WithShutdown(shutdown.ShutdownFunc(consumer.GetConsumer().Stop)),
	)

	return s
}
