package pumpserver

import (
	"context"
	"iam-pump/internal/pumpserver/pump"
	"iam-pump/internal/pumpserver/store"
	"iam-pump/internal/pumpserver/store/mongo"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/che-kwas/iam-kit/server"
)

type pumpServer struct {
	*server.Server
	name     string
	ctx      context.Context
	cancel   context.CancelFunc
	pumpOpts *pump.PumpOptions
	log      *logger.Logger

	err error
}

// NewServer builds a new pumpServer.
func NewServer(name string) *pumpServer {
	ctx, cancel := context.WithCancel(context.Background())

	s := &pumpServer{
		name:     name,
		ctx:      ctx,
		cancel:   cancel,
		pumpOpts: pump.NewPumpOptions(),
		log:      logger.L(),
	}

	return s.initStore().initPump().newServer()
}

// Run runs the pumpServer.
func (s *pumpServer) Run() {
	defer s.cancel()
	defer s.log.Sync()
	if cli := store.Client(); cli != nil {
		defer cli.Close(s.ctx)
	}

	if s.err != nil {
		s.log.Fatal("failed to build the server: ", s.err)
	}

	if err := s.Server.Run(); err != nil {
		s.log.Fatal("server stopped unexpectedly: ", err)
	}
}

func (s *pumpServer) initStore() *pumpServer {
	var storeIns store.Store
	if storeIns, s.err = mongo.MongoStore(s.ctx); s.err != nil {
		return s
	}
	store.SetClient(storeIns)

	return s
}

func (s *pumpServer) initPump() *pumpServer {
	if s.err != nil {
		return s
	}

	go pump.InitPump(s.ctx, s.pumpOpts).Start()
	return s
}

func (s *pumpServer) newServer() *pumpServer {
	if s.err != nil {
		return s
	}

	s.Server, s.err = server.NewServer(s.name)
	return s
}
