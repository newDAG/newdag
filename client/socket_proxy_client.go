package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/newdag/ledger"
)

type SocketProxyClient struct {
	nodeAddr string
	timeout  time.Duration
}

func NewSocketProxyClient(nodeAddr string, timeout time.Duration) *SocketProxyClient {
	return &SocketProxyClient{
		nodeAddr: nodeAddr,
		timeout:  timeout,
	}
}

func (p *SocketProxyClient) getConnection() (*rpc.Client, error) {
	conn, err := net.DialTimeout("tcp", p.nodeAddr, p.timeout)
	if err != nil {
		return nil, err
	}
	return jsonrpc.NewClient(conn), nil
}

func (p *SocketProxyClient) SubmitTx(tx ledger.Transaction) (*bool, error) {
	rpcConn, err := p.getConnection()
	if err != nil {
		return nil, err
	}
	var ack bool
	err = rpcConn.Call("NewDAG.SubmitTx", tx, &ack)
	if err != nil {
		return nil, err
	}
	return &ack, nil
}
