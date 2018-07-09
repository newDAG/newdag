package main

import (
	"fmt"
	"time"

	"github.com/newdag/ledger"
	"github.com/newdag/store"
)

type SocketProxy struct {
	nodeAddress string
	bindAddress string

	client *SocketProxyClient
	server *SocketProxyServer
}

func NewSocketProxy(nodeAddr string, bindAddr string, timeout time.Duration, logger *Logger) (*SocketProxy, error) {
	if logger == nil {
		logger = store.Logger_new("debug")
	}

	client := NewSocketProxyClient(nodeAddr, timeout)
	server, err := NewSocketProxyServer(bindAddr, timeout, logger)
	if err != nil {
		return nil, err
	}

	proxy := &SocketProxy{nodeAddress: nodeAddr, bindAddress: bindAddr, client: client, server: server}
	go proxy.server.listen()

	return proxy, nil
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//Implement ClientProxy interface

func (p *SocketProxy) CommitCh() chan Commit {
	return p.server.commitCh
}

func (p *SocketProxy) SubmitTx(tx ledger.Transaction) error {
	ack, err := p.client.SubmitTx(tx)
	if err != nil {
		return err
	}
	if !*ack {
		return fmt.Errorf("Failed to deliver transaction to Client")
	}
	return nil
}
