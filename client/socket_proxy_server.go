package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/newdag/ledger"
)

// CommitResponse captures both a response and a potential error.
type CommitResponse struct {
	Hash  []byte
	Error error
}

// Commit provides a response mechanism.
type Commit struct {
	Block    ledger.Block
	RespChan chan<- CommitResponse
}

// Respond is used to respond with a response, error or both
func (r *Commit) Respond(stateHash []byte, err error) {
	r.RespChan <- CommitResponse{stateHash, err}
}

type SocketProxyServer struct {
	netListener *net.Listener
	rpcServer   *rpc.Server
	commitCh    chan Commit
	timeout     time.Duration
	logger      *Logger
}

func NewSocketProxyServer(bindAddress string, timeout time.Duration, logger *Logger) (*SocketProxyServer, error) {
	server := &SocketProxyServer{commitCh: make(chan Commit), timeout: timeout, logger: logger}
	if err := server.register(bindAddress); err != nil {
		return nil, err
	}
	return server, nil
}

func (p *SocketProxyServer) register(bindAddress string) error {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterName("State", p)
	p.rpcServer = rpcServer

	l, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return err
	}

	p.netListener = &l

	return nil
}

func (p *SocketProxyServer) listen() error {
	for {
		conn, err := (*p.netListener).Accept()
		if err != nil {
			return err
		}

		go (*p.rpcServer).ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func (p *SocketProxyServer) CommitBlock(block ledger.Block, stateHash *ledger.StateHash) (err error) {
	// Send the Commit over
	respCh := make(chan CommitResponse)
	p.commitCh <- Commit{
		Block:    block,
		RespChan: respCh,
	}

	// Wait for a response
	select {
	case commitResp := <-respCh:
		stateHash.Hash = commitResp.Hash
		if commitResp.Error != nil {
			err = commitResp.Error
		}
	case <-time.After(p.timeout):
		err = fmt.Errorf("command timed out")
	}

	p.logger.Log("debug", "ClientProxyServer.CommitBlock", LogFields{"block": block.Index, "state_hash": stateHash.Hash, "err": err})

	return

}
