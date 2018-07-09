package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	_ "net/http/pprof"

	cli "gopkg.in/urfave/cli.v1"

	hg "github.com/newdag/consensus"
	"github.com/newdag/crypto"
	"github.com/newdag/ledger"
	"github.com/newdag/net"
	"github.com/newdag/node"
	"github.com/newdag/proxy"
	nstore "github.com/newdag/store"
	"github.com/newdag/wallet"
)

type Logger = nstore.Logger
type LogFields = nstore.LogFields

var (
	IPFlag = cli.StringFlag{
		Name:  "ip",
		Usage: "ip",
		Value: "127.0.0.1",
	}
	NodePortFlag = cli.IntFlag{
		Name:  "node_port",
		Usage: "Port to bind node",
		Value: 1337,
	}
	ProxyPortFlag = cli.IntFlag{
		Name:  "proxy_port",
		Usage: "Port to bind Proxy Server",
		Value: 1338,
	}
	ServicePortFlag = cli.IntFlag{
		Name:  "service_port",
		Usage: "Port of HTTP Service",
		Value: 8000,
	}
	ClientAddressFlag = cli.StringFlag{
		Name:  "client_addr",
		Usage: "IP:Port of Client",
		Value: "127.0.0.1:1339",
	}
	GenesisAddrFlag = cli.StringFlag{
		Name:  "address",
		Usage: "Genesis reward address",
		Value: "",
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "NewDAG server"
	app.Usage = "newDAG server with new dag consensus"
	app.HideVersion = true //there is a special command to print the version
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Run node",
			Action: run,
			Flags:  []cli.Flag{IPFlag, NodePortFlag, ProxyPortFlag, ServicePortFlag, ClientAddressFlag},
		},
		{
			Name:   "keygen",
			Usage:  "Dump new key pair",
			Action: keygen,
		},
		{
			Name:   "init",
			Usage:  "Init genesis block",
			Flags:  []cli.Flag{GenesisAddrFlag},
			Action: init_genesis_block,
		},
		{
			Name:  "version",
			Usage: "Show version info",
			Action: func(c *cli.Context) error {
				fmt.Println(nstore.GetVersion())
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	logger := nstore.Logger_new("debug")
	sIP := c.String(IPFlag.Name)
	conf := nstore.NewConfig(c.Int(NodePortFlag.Name), c.Int(ProxyPortFlag.Name), c.Int(ServicePortFlag.Name), c.String(ClientAddressFlag.Name), logger)
	if conf == nil {
		return cli.NewExitError("can't found config file", 1)
	}

	pemKey := crypto.NewPemKey()
	privkey, err := pemKey.ReadKey()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	peerStore := net.NewJSONPeers()
	peers, err := peerStore.Peers()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// There should be at least two peers
	if len(peers) < 2 {
		return cli.NewExitError("peers.json should define at least two peers", 1)
	}

	//Sort peers by public key and assign them an int ID
	//Every participant in the network will run this and assign the same IDs
	sort.Sort(net.ByPubKey(peers))
	pmap := make(map[string]int)
	for i, p := range peers {
		pmap[p.PubKeyHex] = i
	}

	//Find the ID of this node
	nodePub := fmt.Sprintf("0x%X", crypto.FromECDSAPub(&privkey.PublicKey))
	nodeID := pmap[nodePub]

	logger.Log("debug", "PARTICIPANTS", LogFields{"pmap": pmap, "id": nodeID})

	//Instantiate the Store (inmem or badger)
	var store hg.Store
	var needBootstrap bool
	switch conf.StoreType {
	case "inmem":
		store = hg.NewInmemStore(pmap, conf.CacheSize)
	case "badger":
		//If the file already exists, load and bootstrap the store using the file
		if _, err := os.Stat(conf.StorePath); err == nil {
			logger.Debug("loading badger store from existing database")
			store, err = hg.LoadBlotStore(conf.CacheSize, conf.StorePath)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("failed to load BlotStore from existing file: %s", err), 1)
			}
			needBootstrap = true
		} else {
			//Otherwise create a new one
			logger.Debug("creating new badger store from fresh database")
			store, err = hg.NewBlotStore(pmap, conf.CacheSize, conf.StorePath)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("failed to create new BlotStore: %s", err), 1)
			}
		}
	default:
		return cli.NewExitError(fmt.Sprintf("invalid store option: %s", conf.StoreType), 1)
	}
	ledger.Init_ServerWallet(conf, nil /*store.DB*/)

	trans, err := net.NewTCPTransport(fmt.Sprintf("%s:%d", sIP, conf.NodePort), nil, conf.MaxPool, conf.TCPTimeout, logger)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	var prox proxy.AppProxy
	if conf.NoClient {
		prox = proxy.NewInmemAppProxy(logger)
	} else {
		prox = proxy.NewSocketAppProxy(conf.ClientAddress, fmt.Sprintf("%s:%d", sIP, conf.ProxyPort), conf.TCPTimeout, logger)
	}

	node := node.NewNode(conf, nodeID, privkey, peers, store, trans, prox)
	if err := node.Init(needBootstrap); err != nil {
		return cli.NewExitError(fmt.Sprintf("failed to initialize node: %s", err), 1)
	}

	//serviceServer := service.NewService(fmt.Sprintf("%s:%d", sIP, conf.ServicePort), node, logger)
	//go serviceServer.Serve()

	node.Run(true)

	return nil
}

func keygen(c *cli.Context) error {
	pemDump, err := crypto.GeneratePemKey()
	if err != nil {
		fmt.Println("Error generating PemDump")
		os.Exit(2)
	}

	fmt.Println("PublicKey:")
	fmt.Println(pemDump.PublicKey)
	fmt.Println("PrivateKey:")
	fmt.Println(pemDump.PrivateKey)

	return nil
}

func init_genesis_block(c *cli.Context) error {
	logger := nstore.Logger_new("debug")
	conf := nstore.NewConfig(0, 0, 0, "", logger)
	if conf == nil {
		return cli.NewExitError("can't found config file", 1)
	}

	genesisAddr := c.String(GenesisAddrFlag.Name)
	if genesisAddr == "" {
		w := wallet.CreateWalletToFile(conf.Datadir)
		genesisAddr = fmt.Sprintf("%s", w.GetAddress())
	}

	if !wallet.ValidateAddress(genesisAddr) {
		log.Panic("ERROR: Address is not valid")
	}
	db := nstore.OpenKvDatabase(-1)
	if db == nil {
		return cli.NewExitError(fmt.Sprintf("can not open kv database!"), 1)
	}

	err := ledger.InitGenesisBlock(db, genesisAddr)
	if err != nil {
		db.Close()
		return cli.NewExitError(fmt.Sprintf("failed to init genesis block: %s", err), 1)
	}
	db.Close()

	fmt.Println("Genesis address:")
	fmt.Println(genesisAddr)
	fmt.Println("Reward:", ledger.GetBalance(db, genesisAddr), "NG")
	return nil
}
