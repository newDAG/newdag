package proxy

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/newdag/ledger"
	"github.com/newdag/store"
)

type Logger = store.Logger
type LogFields = store.LogFields

//-------client---------------------------------------------------------------------------------------------

type socketAppProxyClient struct {
	clientAddr string
	timeout    time.Duration
	logger     *Logger
}

//call the client rpc State.CommitBlock
func (p *socketAppProxyClient) commitBlock(block ledger.Block) ([]byte, error) {
	var stateHash ledger.StateHash

	conn, err := net.DialTimeout("tcp", p.clientAddr, p.timeout)
	if err != nil {
		return nil, err
	}
	rpcConn := jsonrpc.NewClient(conn)

	err = rpcConn.Call("State.CommitBlock", block, &stateHash)

	p.logger.Log("debug", "AppProxyClient.commitBlock", LogFields{"block": block.Index, "state_hash": stateHash.Hash})

	return stateHash.Hash, err
}

//-------server---------------------------------------------------------------------------------------------

type socketAppProxyServer struct {
	netListener *net.Listener
	rpcServer   *rpc.Server
	submitCh    chan ledger.Transaction
	logger      *Logger
}

func newSocketAppProxyServer(bindAddress string, logger *Logger) *socketAppProxyServer {
	server := &socketAppProxyServer{submitCh: make(chan ledger.Transaction), logger: logger}
	rpcServer := rpc.NewServer()
	rpcServer.RegisterName("NewDAG", server)
	server.rpcServer = rpcServer

	l, err := net.Listen("tcp", bindAddress)
	if err != nil {
		logger.WithField("error", err).Error("Failed to listen")
	}
	server.netListener = &l

	return server
}

//call by rpc
func (p *socketAppProxyServer) SubmitTx(tx ledger.Transaction, ack *bool) error {
	p.logger.Debug("SubmitTx")
	p.submitCh <- tx
	*ack = true
	return nil
}

func (p *socketAppProxyServer) listen() {
	for {
		conn, err := (*p.netListener).Accept()
		if err != nil {
			p.logger.WithField("error", err).Error("Failed to accept")
		}

		go (*p.rpcServer).ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

//----------------------------------------------------------------------------------------------------
type SocketAppProxy struct {
	clientAddress string
	bindAddress   string

	client *socketAppProxyClient
	server *socketAppProxyServer

	logger *Logger
}

func NewSocketAppProxy(clientAddr string, bindAddr string, timeout time.Duration, logger *Logger) *SocketAppProxy {
	if logger == nil {
		logger = store.Logger_new("")
	}

	client := &socketAppProxyClient{clientAddr: clientAddr, timeout: timeout, logger: logger}
	server := newSocketAppProxyServer(bindAddr, logger)

	proxy := &SocketAppProxy{
		clientAddress: clientAddr,
		bindAddress:   bindAddr,
		client:        client,
		server:        server,
		logger:        logger,
	}
	go proxy.server.listen()

	return proxy
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//Implement AppProxy Interface
type AppProxy interface {
	SubmitCh() chan ledger.Transaction
	CommitBlock(block ledger.Block) ([]byte, error)
}

func (p *SocketAppProxy) SubmitCh() chan ledger.Transaction {
	return p.server.submitCh
}

func (p *SocketAppProxy) CommitBlock(block ledger.Block) ([]byte, error) {
	return p.client.commitBlock(block)
}
