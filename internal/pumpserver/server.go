package pumpserver

import (
	"context"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/che-kwas/iam-kit/server"
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

	return s.initStore().initCache().newServer()
}

// Run runs the pumpServer.
func (s *pumpServer) Run() {
	defer s.log.Sync()
	// defer store.Client().Close()
	defer s.cancel()

	if s.err != nil {
		s.log.Fatal("failed to build the server: ", s.err)
	}

	if err := s.Server.Run(); err != nil {
		s.log.Fatal("server stopped unexpectedly: ", err)
	}
}

func (s *pumpServer) initStore() *pumpServer {

	return s
}

func (s *pumpServer) initCache() *pumpServer {
	if s.err != nil {
		return s
	}

	return s
}

func (s *pumpServer) newServer() *pumpServer {
	if s.err != nil {
		return s
	}

	return s
}
