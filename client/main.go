package main

import (
	"fmt"
	"os"
	"time"

	nstore "github.com/newdag/store"

	"github.com/newdag/wallet"
	"gopkg.in/urfave/cli.v1"
)

var (
	ProxyAddressFlag = cli.StringFlag{
		Name:  "proxy_addr",
		Usage: "IP:Port to bind Proxy Server",
		Value: "127.0.0.1:1338",
	}
	ClientPortFlag = cli.IntFlag{
		Name:  "client_port",
		Usage: "IP:Port of Client",
		Value: 1339,
	}
)

type LogFields = nstore.LogFields
type Logger = nstore.Logger

func main() {
	app := cli.NewApp()
	app.Name = "NewDAG client"
	app.Usage = "NewDAG client"
	app.Commands = []cli.Command{
		{
			Name:  "version",
			Usage: "Show version info",
			Action: func(c *cli.Context) error {
				fmt.Println(nstore.GetVersion())
				return nil
			},
		},
		{
			Name:   "wgen",
			Usage:  "create a wallet file",
			Action: wallet_gen,
		},
		{
			Name:   "run",
			Usage:  "Run NewDAG client",
			Flags:  []cli.Flag{ProxyAddressFlag, ClientPortFlag},
			Action: run,
		},
	}

	app.Run(os.Args)
}

func wallet_gen(c *cli.Context) error {
	w := wallet.CreateWalletToFile("")
	fmt.Println("address:", string(w.GetAddress()))
	return nil
}

func run(c *cli.Context) error {
	logger := nstore.Logger_new("debug")

	proxyAddress := c.String(ProxyAddressFlag.Name)
	clientPort := c.Int(ClientPortFlag.Name)
	httpPort := "1388" //c.String(HttpListPortFlag.Name)

	logger.Log("debug", "RUN", LogFields{"proxy_addr": proxyAddress, "client_port": clientPort, "http_port": httpPort})

	if err := Init_ClientWallet(); err != nil {
		return cli.NewExitError(fmt.Sprintf("failed to load BlotStore from existing file: %s", err), 1)
	}

	proxy, err := NewSocketProxy(proxyAddress, fmt.Sprintf(":%d", clientPort), 1*time.Second, logger)
	if err != nil {
		return err
	}
	go block_Run(proxy, logger)

	http_server_main(httpPort, proxy)
	return nil
}

//------------------------------------------------------------------------------------------------------------
func block_Run(proxy *SocketProxy, logger *Logger) {
	for {
		select {
		case commit := <-proxy.CommitCh():
			logger.Log("debug", "CommitBlock", LogFields{"block": commit.Block})
			err := Write_ClientBlock(commit.Block, logger)
			commit.Respond(commit.Block.Hash, err)
		}
	}
}
